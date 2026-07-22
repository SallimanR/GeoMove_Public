-- name: CreateDriver :exec
INSERT INTO
	driver (
		user_id,
		name,
		phone,
		work_starts,
		work_ends,
		location,
		rating,
		address
	) VALUES ($1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint(@lon::REAL, @lat::REAL), 4326), $6, $7);

-- name: UpdateDriverProfileImage :exec
UPDATE driver
SET profile_image = $2, updated_at = NOW()
WHERE user_id = $1;

-- name: GetDriverByUserID :one
SELECT
	d.user_id,
	d.name,
	d.phone,
	d.profile_image,
	d.work_starts,
	d.work_ends,
	d.is_available,
	d.last_seen,
	d.rating,
	ST_X(d.location::geometry)::REAL as lon,
	ST_Y(d.location::geometry)::REAL as lat,
	COALESCE(t.max_car_weight_kg, 0) as max_car_weight_kg,
	COALESCE(t.max_car_length_meters, 0) as max_car_length_meters,
	COALESCE(d.address, '') as address
FROM driver d
LEFT JOIN tow_driver t ON t.driver_id = d.user_id
WHERE d.user_id = $1;

-- name: GetFilteredDrivers :many
SELECT 
    user_id,
	name,
	phone,
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

-- name: CreateTowDriver :exec
INSERT INTO tow_driver (driver_id, max_car_weight_kg, max_car_length_meters)
VALUES ($1, $2, $3);

-- name: UpdateDriver :exec
UPDATE driver
SET name = $2, phone = $3, work_starts = $4, work_ends = $5,
    location = ST_SetSRID(ST_MakePoint($6::REAL, $7::REAL), 4326),
    address = $8,
    updated_at = NOW()
WHERE user_id = $1;

-- name: UpsertTowDriver :exec
INSERT INTO tow_driver (driver_id, max_car_weight_kg, max_car_length_meters)
VALUES ($1, $2, $3)
ON CONFLICT (driver_id) DO UPDATE SET
	max_car_weight_kg = EXCLUDED.max_car_weight_kg,
	max_car_length_meters = EXCLUDED.max_car_length_meters;
