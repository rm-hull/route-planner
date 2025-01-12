package repository

import (
	"context"
	"fmt"
	"regexp"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rm-hull/route-planner/models"
)

type RefDataRepository interface {
	Store(ctx context.Context, refData *models.RefData) error
	FetchAll(ctx context.Context) (*map[string]models.RefData, error)
}

type RefDataRepositoryImpl struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewRefDataRepository(pool *pgxpool.Pool, tableName string) *RefDataRepositoryImpl {
	return &RefDataRepositoryImpl{pool: pool, tableName: tableName}
}

func (repo *RefDataRepositoryImpl) Store(ctx context.Context, refData *models.RefData) error {
	sql := fmt.Sprintf(`
		INSERT INTO "%s" (value, description) VALUES ($1, $2)
		ON CONFLICT (value) DO UPDATE SET description = EXCLUDED.description
	`, repo.tableName)

	_, err := repo.pool.Exec(ctx, sql, refData.Value, nullableString(*refData.Description))
	return err
}

var re *regexp.Regexp = regexp.MustCompile(`\s+`)

func nullableString(text string) *string {
	var nullableText *string = nil
	if re.ReplaceAllString(text, "") != "" {
		nullableText = &text
	}
	return nullableText
}

func (repo *RefDataRepositoryImpl) FetchAll(ctx context.Context) (*map[string]models.RefData, error) {
	sql := fmt.Sprintf(`SELECT id, value, description FROM "%s"`, repo.tableName)

	rows, err := repo.pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ref data: %v", err)
	}
	defer rows.Close()

	results := make(map[string]models.RefData, 0)
	for rows.Next() {
		var refData models.RefData
		if err := rows.Scan(&refData.ID, &refData.Value, &refData.Description); err != nil {
			return nil, fmt.Errorf("failed to scan ref data: %v", err)
		}
		results[refData.Value] = refData
	}

	return &results, nil
}
