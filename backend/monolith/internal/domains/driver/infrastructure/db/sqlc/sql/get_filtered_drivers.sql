-- name: GetFilteredDrivers :many
SELECT 
    id,
	name,
    work_starts, 
    work_ends, 
    rating, 
	ST_X(location::GEOMETRY)::REAL AS lon,
	ST_Y(location::GEOMETRY)::REAL AS lat,
    st_distance(
        location, 
        st_setsrid(
            st_makepoint(@lat::REAL, @lon::REAL), 
            4326
        )::geometry
    )::real AS distance
FROM driver
WHERE 
    (sqlc.narg('work_starts')::TIME IS NULL OR work_starts <= sqlc.narg('work_starts'))
    AND (sqlc.narg('work_ends')::TIME IS NULL OR work_ends >= sqlc.narg('work_ends'))
    AND (sqlc.narg('min_rating')::REAL IS NULL OR rating >= sqlc.narg('min_rating'))
ORDER BY distance;
