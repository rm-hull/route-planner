# route-planner


## Database Schema Creation

```sql
CREATE SCHEMA IF NOT EXISTS route_planner;

GRANT USAGE ON SCHEMA route_planner TO homelab;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA route_planner TO homelab;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA route_planner TO homelab;

ALTER DEFAULT PRIVILEGES IN SCHEMA route_planner GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO homelab;
ALTER DEFAULT PRIVILEGES IN SCHEMA route_planner GRANT USAGE, SELECT ON SEQUENCES TO homelab;
```

# Migration and ref-data import

```bash
route-planner migration db/migrations
route-planner refdata road_classifications http://www.os.uk/xml/codelists/RoadClassificationValue.xml
route-planner refdata road_functions http://www.os.uk/xml/codelists/RoadFunctionValue.xml
route-planner refdata form_of_way_types http://www.os.uk/xml/codelists/FormOfWayTypeValue.xml
route-planner refdata form_of_road_types https://raw.githubusercontent.com/rm-hull/route-planner/refs/heads/main/data/FormOfRoadNodeTypeValue.xml
