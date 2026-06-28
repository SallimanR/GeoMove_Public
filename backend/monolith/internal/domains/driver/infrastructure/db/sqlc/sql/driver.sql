-- name: CreateDriver :one
INSERT INTO
	driver (
		work_starts,
		work_ends,
		location,
		rating
	) VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3::REAL, $4::REAL), 4326), $5)
RETURNING id;

-- name: GetDriverByID :one
SELECT
	id,
	work_starts,
	work_ends,
	rating,
	ST_X(location::geometry)::REAL as lon,
	ST_Y(location::geometry)::REAL as lat
FROM driver
WHERE id = $1;
