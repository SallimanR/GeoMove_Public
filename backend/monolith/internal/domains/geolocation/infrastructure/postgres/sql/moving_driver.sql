-- name: UpdateMovingDriver :exec
INSERT INTO moving_driver (
    driver_id,
    updated_at,
	realtime_location,
	travel_time,
	path_meters
) VALUES (
	$1,
    now(),
	ST_SetSRID(ST_MakePoint(@lon::REAL, @lat::REAL), 4326),
    $2,
	$3
)
ON CONFLICT (driver_id) DO UPDATE SET
    updated_at = EXCLUDED.updated_at,
    realtime_location = EXCLUDED.realtime_location,
    travel_time = EXCLUDED.travel_time,
	path_meters = EXCLUDED.path_meters;

-- name: GetMovingDriverByID :one
SELECT 
	driver_id,
	ST_X(realtime_location::GEOMETRY)::REAL AS lon,
	ST_Y(realtime_location::GEOMETRY)::REAL AS lat,
	travel_time,
	path_meters
FROM moving_driver
WHERE driver_id = $1;

-- name: GetClosestWithinRadiusMovingDriver :many
WITH user_location AS (
	SELECT st_setsrid(st_makepoint(@lon::REAL, @lat::REAL), 4326)::GEOGRAPHY AS geog
)
SELECT
	md.driver_id,
	st_distance(md.realtime_location, user_location.geog)::INTEGER AS distance_meters,
	ST_X(md.realtime_location::GEOMETRY)::REAL AS lon,
	ST_Y(md.realtime_location::GEOMETRY)::REAL AS lat,
	travel_time,
	path_meters
FROM moving_driver AS md
CROSS JOIN user_location
WHERE st_dwithin(md.realtime_location, user_location.geog, @radius::INTEGER)
ORDER BY md.realtime_location <-> user_location.geog
LIMIT 20;
