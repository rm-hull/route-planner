CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS pgrouting;

CREATE TABLE road_classifications (
    id SERIAL PRIMARY KEY,
    value TEXT NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE road_functions (
    id SERIAL PRIMARY KEY,
    value TEXT NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE form_of_way_types (
    id SERIAL PRIMARY KEY,
    value TEXT NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE form_of_road_types (
    id SERIAL PRIMARY KEY,
    value TEXT NOT NULL UNIQUE,
    description TEXT
);
