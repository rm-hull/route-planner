package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rm-hull/route-planner/models"
	"github.com/spaolacci/murmur3"
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
	if _, err := repo.pool.Exec(ctx, `ALTER TABLE road_links DISABLE TRIGGER ALL`); err != nil {
		return err
	}
	if _, err := repo.pool.Exec(ctx, `ALTER TABLE road_nodes DISABLE TRIGGER ALL`); err != nil {
		return err
	}
	return nil
}

func (repo *GmlRepositoryImpl) EnableTriggers(ctx context.Context) error {
	if _, err := repo.pool.Exec(ctx, `ALTER TABLE road_links ENABLE TRIGGER ALL`); err != nil {
		return err
	}
	if _, err := repo.pool.Exec(ctx, `ALTER TABLE road_nodes ENABLE TRIGGER ALL`); err != nil {
		return err
	}
	return nil
}

func hash(id string) int64 {
	return int64(murmur3.Sum64([]byte(id[2:])))
}

func (repo *GmlRepositoryImpl) StoreRoadNodes(ctx context.Context, roadNodes ...models.RoadNode) error {
	sql := `
		INSERT INTO road_nodes (id, gml_id, location, form_of_road_id)
		VALUES ($1, $2, ST_Transform(ST_SetSRID(ST_GeomFromText($3), 27700), 4326), $4)
		ON CONFLICT (id) DO UPDATE SET
			location = EXCLUDED.location, form_of_road_id = EXCLUDED.form_of_road_id;
	`

	batch := &pgx.Batch{}

	for _, roadNode := range roadNodes {
		batch.Queue(sql,
			hash(roadNode.ID),
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
			return fmt.Errorf("batch insert on road_nodes failed at query %d (gml:id=%s): %v", i, roadNodes[i].ID, err)
		}
	}

	return nil
}

func (repo *GmlRepositoryImpl) StoreRoadLinks(ctx context.Context, roadLinks ...models.RoadLink) error {
	sql := `
		INSERT INTO road_links (
			id, source_id, target_id, gml_id, center_line, start_node_id, end_node_id, road_classification_id, road_function_id,
			form_of_way_id, road_classification_number, name1, length_m, loop, primary_route, trunk_road)
		VALUES (
			$1, $2, $3, $4, ST_Transform(ST_SetSRID(ST_GeomFromText($5), 27700), 4326), $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16
		)
		ON CONFLICT (id) DO UPDATE SET
			source_id = EXCLUDED.source_id, target_id = EXCLUDED.target_id, gml_id = EXCLUDED.gml_id,
			center_line = EXCLUDED.center_line, start_node_id = EXCLUDED.start_node_id, end_node_id = EXCLUDED.end_node_id,
			road_classification_id = EXCLUDED.road_classification_id, road_function_id = EXCLUDED.road_function_id,
			form_of_way_id = EXCLUDED.form_of_way_id, road_classification_number = EXCLUDED.road_classification_number,
			name1 = EXCLUDED.name1, length_m = EXCLUDED.length_m, loop = EXCLUDED.loop, primary_route = EXCLUDED.primary_route,
			trunk_road = EXCLUDED.trunk_road;
	`

	batch := &pgx.Batch{}

	for _, roadLink := range roadLinks {

		batch.Queue(sql,
			hash(roadLink.ID),
			hash(roadLink.StartNode.Ref()),
			hash(roadLink.EndNode.Ref()),
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
			return fmt.Errorf("batch insert on road_links failed at query %d (gml:id=%s): %v", i, roadLinks[i].ID, err)
		}
	}

	return nil
}

func (repo *GmlRepositoryImpl) StoreMotorwayJunctions(ctx context.Context, motorwayJunctions ...models.MotorwayJunction) error {
	// panic("not implemented")
	return nil
}
