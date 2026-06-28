-- name: GetFilteredDriversAndTiles :many
WITH
	tile_bounds AS (SELECT st_tileenvelope(:z, :x, :y) AS geom_3857),
	filtered_drivers AS (
		SELECT 
			work_starts, 
			work_ends, 
			rating, 
			location,
			st_distance(location, 'SRID=4326;POINT(37.62639302057005 55.74330433694732)'::geometry) as distance
		FROM driver
		WHERE 
			($1::time IS NULL OR work_starts >= $1) AND
			($2::time IS NULL OR work_ends <= $2) AND
			($3::float IS NULL OR rating >= $3)
	),
	tile_data AS (
		SELECT
			ST_AsMVTGeom(
				ST_Transform(ST_CurveToLine(location::geometry), 3857),
				tile_bounds.geom_3857,
				4096, 64, true
			) as geom,
			work_starts,
			work_ends,
			rating,
			distance
			-- row_to_json(fd) as driver_data  -- Store complete driver data
		FROM filtered_drivers, tile_bounds
		WHERE location && ST_Transform(tile_bounds.geom_3857, 4326)
	)
SELECT 
	(
		SELECT json_agg(
			json_build_object(
				'work_starts', work_starts,
				'work_ends', work_ends, 
				'rating', rating,
				'location', ST_AsGeoJSON(location)::json,
				'distance', distance
			)
			ORDER BY distance
		) 
		FROM filtered_drivers
	) AS filtered_data_json,
	
	ST_AsMVT(tile_data, 'driver', 4096, 'geom') AS mvt_tile
	
	-- (SELECT gzip(ST_AsMVT(tile_data, 'drivers', 4096, 'geom'))) as mvt_tile_gzipped
FROM tile_data
WHERE geom IS NOT NULL;
