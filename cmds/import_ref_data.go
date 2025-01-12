package cmds

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/rm-hull/route-planner/db"
	"github.com/rm-hull/route-planner/models"
)

func ImportRefData(tableName string, url string) error {

	dict, err := parse(url)
	if err != nil {
		return fmt.Errorf("failed to parse ref data: %v", err)
	}

	err = insertIntoDb(tableName, dict)
	if err != nil {
		return fmt.Errorf("failed to insert ref data: %v", err)
	}

	log.Printf("Imported %d records into table: %s\n", len(dict.Entries), tableName)
	return nil
}

func parse(url string) (*models.Dictionary, error) {

	xmlData, err := downloadFile(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file '%s': %v", url, err)
	}

	var dict models.Dictionary
	err = xml.Unmarshal(xmlData, &dict)
	if err != nil {
		return nil, err
	}

	return &dict, nil
}

func downloadFile(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Route Planner (https://github.com/rm-hull/route-planner)")
	req.Header.Set("Accept", "application/xml, text/xml")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}

func insertIntoDb(tableName string, dict *models.Dictionary) error {
	config := db.ConfigFromEnv()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	pool, err := db.NewDBPool(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	sql := fmt.Sprintf(`
		INSERT INTO "%s"."%s" (value, description) VALUES ($1, $2)
		ON CONFLICT (value) DO UPDATE SET description = EXCLUDED.description
	`, config.Schema, tableName)

	batch := &pgx.Batch{}
	re := regexp.MustCompile(`\s+`)

	for _, entry := range dict.Entries {
		var description *string = nil
		if re.ReplaceAllString(entry.Definition.Description, "") != "" {
			description = &entry.Definition.Description
		}
		batch.Queue(sql, entry.Definition.Identifier.Value, description)
	}

	results := pool.SendBatch(ctx, batch)
	defer results.Close()

	// Ensure all queries in the batch succeeded
	for i := range batch.Len() {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("batch insert failed at query %d: %v", i, err)
		}
	}

	return nil
}
