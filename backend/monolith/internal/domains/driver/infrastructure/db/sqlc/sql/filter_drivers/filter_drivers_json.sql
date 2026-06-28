-- name: GetFilteredDriversJson :many
WITH
	filtered_drivers AS (
		SELECT 
			work_starts, 
			work_ends, 
			rating, 
			location,
			st_distance(location, 'SRID=4326;POINT(37.62639302057005 55.74330433694732)'::geometry) AS distance
		FROM driver
		WHERE 
			($1::time IS NULL OR work_starts >= $1) AND
			($2::time IS NULL OR work_ends <= $2) AND
			($3::float IS NULL OR rating >= $3)
	)
SELECT json_agg(
	json_build_object(
		'work_starts', work_starts,
		'work_ends', work_ends, 
		'rating', rating,
		'location', ST_AsGeoJSON(location)::json,
		'distance', distance
	)
	ORDER BY distance
) AS drivers
FROM filtered_drivers;
