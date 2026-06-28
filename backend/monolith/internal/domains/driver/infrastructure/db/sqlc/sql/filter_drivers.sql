-- name: GetFilteredDrivers :many
SELECT 
	id,
	work_starts, 
	work_ends, 
	rating, 
	location::text,
	st_distance(location, st_setsrid(st_makepoint($1::real, $2::real), 4326)::geometry)::real AS distance
FROM driver
WHERE 
	($3::text IS NULL OR work_starts >= $3::time) AND
	($4::text IS NULL OR work_ends <= $4::time) AND
	($5::real IS NULL OR rating >= $5)
ORDER BY distance;
