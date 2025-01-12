package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rm-hull/route-planner/models"
)

type GmlRepository interface {
	StoreRoadNodes(ctx context.Context, roadNodes ...models.RoadNode) error
	StoreRoadLinks(ctx context.Context, roadLinks ...models.RoadLink) error
	StoreMotorwayJunctions(ctx context.Context, junctions ...models.MotorwayJunction) error
}

type GmlRepositoryImpl struct {
	pool                *pgxpool.Pool
	roadClassifications map[string]models.RefData
	roadFunctions       map[string]models.RefData
	formOfWayTypes      map[string]models.RefData
	formOfRoadTypes     map[string]models.RefData
}

func NewGmlRepository(pool *pgxpool.Pool) (*GmlRepositoryImpl, error) {
	ctx := context.Background()
	roadClassifications, err := NewRefDataRepository(pool, "road_classifications").FetchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching road_classifications: %v", err)
	}

	roadFunctions, err := NewRefDataRepository(pool, "road_functions").FetchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching road_functions: %v", err)
	}

	formOfWayTypes, err := NewRefDataRepository(pool, "form_of_way_types").FetchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching form_of_way_types: %v", err)
	}

	formOfRoadTypes, err := NewRefDataRepository(pool, "form_of_road_types").FetchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching form_of_road_types: %v", err)
	}

	return &GmlRepositoryImpl{
		pool:                pool,
		roadClassifications: *roadClassifications,
		roadFunctions:       *roadFunctions,
		formOfWayTypes:      *formOfWayTypes,
		formOfRoadTypes:     *formOfRoadTypes,
	}, nil
}

func (repo *GmlRepositoryImpl) DisableTriggers(ctx context.Context) error {
	_, err := repo.pool.Exec(ctx, `ALTER TABLE road_links DISABLE TRIGGER ALL`)
	return err
}

func (repo *GmlRepositoryImpl) EnableTriggers(ctx context.Context) error {
	_, err := repo.pool.Exec(ctx, `ALTER TABLE road_links ENABLE TRIGGER ALL`)
	return err
}

func (repo *GmlRepositoryImpl) StoreRoadNodes(ctx context.Context, roadNodes ...models.RoadNode) error {
	sql := `
		INSERT INTO road_nodes (id, location, form_of_road_id)
		VALUES ($1, ST_Transform(ST_SetSRID(ST_GeomFromText($2), 27700), 4326), $3)
		ON CONFLICT (id) DO UPDATE SET
			location = EXCLUDED.location, form_of_road_id = EXCLUDED.form_of_road_id;
	`

	batch := &pgx.Batch{}

	for _, roadNode := range roadNodes {
		batch.Queue(sql,
			roadNode.ID,
			roadNode.Geometry.AsPoint(),
			repo.formOfRoadTypes[roadNode.FormOfRoadNode.Value].ID,
		)
	}

	results := repo.pool.SendBatch(ctx, batch)
	defer results.Close()

	// Ensure all queries in the batch succeed
	for i := range batch.Len() {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("batch insert failed at query %d: %v", i, err)
		}
	}

	return nil
}

func (repo *GmlRepositoryImpl) StoreRoadLinks(ctx context.Context, roadLinks ...models.RoadLink) error {
	sql := `
		INSERT INTO road_links (
			id, center_line, start_node_id, end_node_id, road_classification_id, road_function_id,
			form_of_way_id, road_classification_number, name1, length_m, loop, primary_route, trunk_road)
		VALUES (
			$1, ST_Transform(ST_SetSRID(ST_GeomFromText($2), 27700), 4326), $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13
		)
		ON CONFLICT (id) DO UPDATE SET
			center_line = EXCLUDED.center_line, start_node_id = EXCLUDED.start_node_id, end_node_id = EXCLUDED.end_node_id,
			road_classification_id = EXCLUDED.road_classification_id, road_function_id = EXCLUDED.road_function_id,
			form_of_way_id = EXCLUDED.form_of_way_id, road_classification_number = EXCLUDED.road_classification_number,
			name1 = EXCLUDED.name1, length_m = EXCLUDED.length_m, loop = EXCLUDED.loop, primary_route = EXCLUDED.primary_route,
			trunk_road = EXCLUDED.trunk_road;
	`

	batch := &pgx.Batch{}

	for _, roadLink := range roadLinks {
		batch.Queue(sql,
			roadLink.ID,
			roadLink.CentrelineGeometry.AsLineString(),
			roadLink.StartNode.Ref(),
			roadLink.EndNode.Ref(),
			repo.roadClassifications[roadLink.RoadClassification.Value].ID,
			repo.roadFunctions[roadLink.RoadFunction.Value].ID,
			repo.formOfWayTypes[roadLink.FormOfWay.Value].ID,
			roadLink.RoadClassificationNumber,
			roadLink.Name1,
			roadLink.Length.ConvertTo("m"),
			roadLink.Loop,
			roadLink.PrimaryRoute,
			roadLink.TrunkRoad,
		)
	}

	results := repo.pool.SendBatch(ctx, batch)
	defer results.Close()

	// Ensure all queries in the batch succeed
	for i := range batch.Len() {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("batch insert failed at query %d: %v", i, err)
		}
	}

	return nil
}

func (repo *GmlRepositoryImpl) StoreMotorwayJunctions(ctx context.Context, motorwayJunctions ...models.MotorwayJunction) error {
	// panic("not implemented")
	return nil
}
