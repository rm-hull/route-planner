package cmds

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rm-hull/route-planner/db"
	"github.com/rm-hull/route-planner/models"
	"github.com/rm-hull/route-planner/repository"
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

	repo := repository.NewRefDataRepository(pool, tableName)
	for _, entry := range dict.Entries {
		err := repo.Store(ctx, &models.RefData{Value: entry.Definition.Identifier.Value, Description: &entry.Definition.Description})
		if err != nil {
			return fmt.Errorf("failed to store record: %v", err)
		}
	}

	return nil
}
