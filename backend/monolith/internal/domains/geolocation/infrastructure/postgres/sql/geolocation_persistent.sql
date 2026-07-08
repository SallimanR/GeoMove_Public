-- name: FindWithinRadiusDrivers :many
WITH
	center AS (
		SELECT 
			$1::REAL AS lat, 
			$2::REAL AS lon,
			$3::INTEGER AS radius,
			ST_SetSRID(ST_MakePoint(lon, lat), 4326)::GEOGRAPHY AS geog,
			-- TODO: change resolution to bigger of h3
			-- h3_latlng_to_cell(lat, lon, 9) AS cell
			h3_latlng_to_cell(geog, 9) AS cell
	),
	covering AS (
		SELECT h3_grid_disk(cell, (ceil(radius / 350.0)::INTEGER + 1)) AS cell
		FROM center
	)
SELECT
	d.id,
	ST_X(d.location::GEOMETRY)::REAL AS lon,
	ST_Y(d.location::GEOMETRY)::REAL AS lat,
	ST_Distance(d.location, center.geog) AS distance_meters
FROM
	driver AS d
	CROSS JOIN center
	JOIN covering ON d.h3_res9 = covering.cell
WHERE
	d.is_available
	AND ST_DWithin(d.location, center.geog, center.radius)
ORDER BY d.location <-> center.geog
LIMIT 50;
