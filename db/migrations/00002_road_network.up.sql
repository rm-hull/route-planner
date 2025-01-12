CREATE TABLE road_nodes (
    id TEXT PRIMARY KEY,
    location GEOMETRY(POINT, 27700) NOT NULL, -- BNG SRID
    form_of_road_id INT NOT NULL REFERENCES form_of_road_types(id)
);

CREATE INDEX idx_road_nodes_geolocation ON road_nodes USING GIST (location);

CREATE TABLE road_links (
    id TEXT PRIMARY KEY,
    center_line GEOMETRY(LINESTRING, 27700) NOT NULL, -- BNG SRID
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    start_node_id TEXT NOT NULL REFERENCES road_nodes(id),
    end_node_id TEXT NOT NULL REFERENCES road_nodes(id),
    road_classification_id INT NOT NULL REFERENCES road_classifications(id),
    road_function_id INT NOT NULL REFERENCES road_functions(id),
    form_of_way_id INT NOT NULL REFERENCES form_of_way_types(id),
    road_classification_number TEXT,
    name1 TEXT,
    length_m NUMERIC(8,2),
    loop BOOLEAN,
    primary_route BOOLEAN,
    trunk_road BOOLEAN,
    road_name_toid TEXT,
    road_number_toid TEXT
);

CREATE INDEX idx_road_links_primary_route ON road_links (primary_route);
CREATE INDEX idx_road_links_trunk_road ON road_links (trunk_road);
CREATE INDEX idx_road_links_center_line ON road_links USING GIST (center_line);



