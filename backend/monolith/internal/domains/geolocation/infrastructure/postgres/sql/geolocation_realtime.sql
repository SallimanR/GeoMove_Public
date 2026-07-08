-----------------------  *Single driver* --------------------

-- name: CreateDriverRealtime :exec
INSERT INTO driver_realtime (driver_id)
VALUES ($1);

-- name: UpdateDriverRealtime :exec
UPDATE driver_realtime
SET
	updated_at = CURRENT_TIMESTAMP,
	realtime_location = st_setsrid(st_makepoint($3::REAL, $2::REAL), 4326), -- lon ($3), lat ($2) order
	predicted_bearing = $4::REAL,
	average_speed = $5::REAL
WHERE
	driver_id = $1;

-- name: GetDriverRealtimeByID :one
SELECT 
	driver_id,
	ST_X(realtime_location::GEOMETRY)::REAL AS lon,
	ST_Y(realtime_location::GEOMETRY)::REAL AS lat
FROM driver_realtime
WHERE 
	driver_id = $1
	-- TODO: do Garbage Collecting instead (Vacuum)
	-- AND realtime_updated_at > NOW() - INTERVAL '5 minutes'
;

-- TODO: replace "CreateRealtime" with this
-- name: CreateOrUpdateDriverRealtime :exec
INSERT INTO driver_realtime (
    driver_id,
	updated_at,
    realtime_location,
    average_speed,
    predicted_bearing
) VALUES (
    $1::INTEGER,
	now(),
    $2::GEOGRAPHY,
    $3::REAL,
    $4::REAL
)
ON CONFLICT (driver_id) DO UPDATE SET
	updated_at = EXCLUDED.updated_at,
    realtime_location = EXCLUDED.realtime_location,
    average_speed = EXCLUDED.average_speed,
    predicted_bearing = EXCLUDED.predicted_bearing;

-----------------------  *Find Within Radius drivers* implementations --------------------

-- name: FindWithinRadiusDriversRealtimeH3 :many
WITH
	center AS (
		SELECT 
			$1::REAL AS lat, 
			$2::REAL AS lon,
			$3::INTEGER AS radius,
			-- TODO: move to params
			$4::INTEGER AS result_limit,
			ST_SetSRID(ST_MakePoint(lon, lat), 4326)::GEOGRAPHY AS geog,
			h3_latlng_to_cell(lat, lon, 9) AS cell
	),
	covering AS (
		SELECT h3_grid_disk(cell, (ceil(radius / 350.0)::INTEGER + 1)) AS cell
		FROM center
	)
SELECT
	dr.driver_id,
	ST_Distance(dr.realtime_location, center.geog) AS distance_meters
FROM driver_realtime AS dr
CROSS JOIN center
JOIN covering ON dr.coarse_h3 = covering.cell
WHERE ST_DWithin(dr.realtime_location, center.geog, center.radius)
ORDER BY dr.realtime_location <-> center.geog
LIMIT center.result_limit;

-----------------------  *Find Closest drivers* implementations --------------------

-- 1. GiST index -> 2. ORDER BY closest "<->"
-- name: FindClosestDriversRealtime :many
WITH point AS (
    SELECT ST_SetSRID(ST_MakePoint($2::REAL, $1::REAL), 4326)::GEOGRAPHY AS geog -- lon ($2), lat ($1) order
)
SELECT
	dr.driver_id,
	ST_Distance(dr.realtime_location, point.geog)::INTEGER AS distance_meters
FROM
	driver_realtime AS dr
	CROSS JOIN point
ORDER BY dr.realtime_location <-> point.geog
LIMIT 10;

-- 1. H3 cells = [area] -> 2. GiST index = [exact distance] withith [area] -> 3. ORDER BY closest "<->"
-- name: FindClosestDriversRealtimeH3 :many
WITH 
    point AS (
        SELECT 
            $1::REAL AS lat,
            $2::REAL AS lon,
            ST_SetSRID(ST_MakePoint($2, $1), 4326)::GEOGRAPHY AS geog
    ),
    coarse_cells AS (
		-- TODO: change resolution to bigger of h3
        SELECT h3_grid_disk(h3_latlng_to_cell(point.geog::GEOMETRY, 2), 1) AS cell
        FROM point
    )
SELECT
	dr.driver_id,
  	ST_Distance(dr.realtime_location, point.geog)::INTEGER AS distance_meters
FROM
	driver_realtime AS dr
	CROSS JOIN point
    JOIN coarse_cells AS cc ON dr.coarse_h3 = cc.cell
ORDER BY dr.realtime_location <-> point.geog
LIMIT 10;

-- name: FindClosestWithinRadiusDriversRealtime :many
WITH user_location AS (
	SELECT st_setsrid(st_makepoint($2::REAL, $1::REAL), 4326)::GEOGRAPHY AS geog -- lon ($2), lat ($1) order
)
SELECT
	dr.driver_id,
	-- TODO:
	-- ST_X(dr.realtime_location::GEOMETRY)::REAL AS lon,
	-- ST_Y(dr.realtime_location::GEOMETRY)::REAL AS lat,
	st_distance(dr.realtime_location, user_location.geog)::INTEGER AS distance_meters
FROM
	driver_realtime AS dr
	CROSS JOIN user_location
WHERE st_dwithin(dr.realtime_location, user_location.geog, $3::INTEGER) -- radius in meters
-- "<->" uses GiST indexing, it is faster over:
-- "ORDER BY distance_meters ASC"
ORDER BY dr.realtime_location <-> user_location.geog
LIMIT 20;
