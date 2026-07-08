-- name: CreateDriver :exec
INSERT INTO
	driver (
		user_id,
		name,
		work_starts,
		work_ends,
		location,
		rating
	) VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint(@lon::REAL, @lat::REAL), 4326), $5);

-- name: GetDriverByUserID :one
SELECT
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

-- name: GetFilteredDrivers :many
SELECT 
    user_id,
	name,
	profile_image,
    work_starts, 
    work_ends, 
	is_available,
	last_seen,
    rating, 
	ST_X(location::GEOMETRY)::REAL AS lon,
	ST_Y(location::GEOMETRY)::REAL AS lat,
    st_distance(
        location, 
        st_setsrid(
            st_makepoint(@lon::REAL, @lat::REAL), 
            4326
        )::geometry
    )::real AS distance
FROM driver
WHERE 
    (sqlc.narg('work_starts')::TIME IS NULL OR work_starts <= sqlc.narg('work_starts'))
    AND (sqlc.narg('work_ends')::TIME IS NULL OR work_ends >= sqlc.narg('work_ends'))
    AND (sqlc.narg('min_rating')::REAL IS NULL OR rating >= sqlc.narg('min_rating'))
ORDER BY distance;
