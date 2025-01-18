CREATE TABLE road_nodes (
    id BIGINT PRIMARY KEY,
    gml_id TEXT NOT NULL UNIQUE,
    location GEOMETRY(POINT, 4326) NOT NULL, -- WSG84 SRID
    form_of_road_id INT NOT NULL REFERENCES form_of_road_types(id)
);

CREATE INDEX idx_road_nodes_geolocation ON road_nodes USING GIST (location);

CREATE TABLE road_links (
    id BIGINT PRIMARY KEY,
    source_id BIGINT NOT NULL REFERENCES road_nodes(id),
    target_id BIGINT NOT NULL REFERENCES road_nodes(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gml_id TEXT NOT NULL UNIQUE,
    start_node_id TEXT NOT NULL REFERENCES road_nodes(gml_id),
    end_node_id TEXT NOT NULL REFERENCES road_nodes(gml_id),
    road_classification_id INT NOT NULL REFERENCES road_classifications(id),
    road_function_id INT NOT NULL REFERENCES road_functions(id),
    form_of_way_id INT NOT NULL REFERENCES form_of_way_types(id),
    road_classification_number TEXT,
    name1 TEXT,
    length_m NUMERIC(8,2),
    loop BOOLEAN,
    primary_route BOOLEAN,
    trunk_road BOOLEAN,
    center_line GEOMETRY(LINESTRING, 4326) NOT NULL -- WSG84 SRID
);

CREATE INDEX idx_road_links_primary_route ON road_links (primary_route);
CREATE INDEX idx_road_links_trunk_road ON road_links (trunk_road);
CREATE INDEX idx_road_links_center_line ON road_links USING GIST (center_line);



