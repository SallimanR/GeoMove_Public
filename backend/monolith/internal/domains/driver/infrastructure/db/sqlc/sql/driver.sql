-- name: CreateDriver :one
INSERT INTO
	driver (
		user_id,
		name,
		work_starts,
		work_ends,
		location,
		rating
	) VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint(@lon::REAL, @lat::REAL), 4326), $5)
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

-- name: GetDriverByUserID :one
SELECT
	id,
	user_id,
	name,
	profile_image,
	work_starts,
	work_ends,
	is_available,
	last_seen,
	rating,
	ST_X(location::geometry)::REAL as lon,
	ST_Y(location::geometry)::REAL as lat
FROM driver
WHERE user_id = $1;
