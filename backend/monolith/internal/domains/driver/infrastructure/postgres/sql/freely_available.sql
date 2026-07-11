-- name: CreateFreelyAvailable :exec
INSERT INTO tow_driver_freely_available (
	user_id,
	from_date,
	to_date,
	from_location,
	en_route_order,
	tariff_per_km,
	from_address
) VALUES (
	$1, $2, $3,
	ST_SetSRID(ST_MakePoint(@from_lon::REAL, @from_lat::REAL), 4326),
	$4, $5, $6
);

-- name: CreateFreelyAvailableLocation :exec
INSERT INTO tow_driver_freely_available_to_location_list (
	tow_driver,
	location,
	address
) VALUES (
	$1,
	ST_SetSRID(ST_MakePoint(@lon::REAL, @lat::REAL), 4326),
	$2
);

-- name: UpdateFreelyAvailable :exec
UPDATE tow_driver_freely_available
SET
	from_date = $2,
	to_date = $3,
	from_location = ST_SetSRID(ST_MakePoint(@from_lon::REAL, @from_lat::REAL), 4326),
	en_route_order = $4,
	tariff_per_km = $5,
	from_address = $6
WHERE user_id = $1;

-- name: DeleteFreelyAvailableLocations :exec
DELETE FROM tow_driver_freely_available_to_location_list
WHERE tow_driver = $1;

-- name: DeleteFreelyAvailable :exec
DELETE FROM tow_driver_freely_available
WHERE user_id = $1;

-- name: GetFreelyAvailableByUserID :one
SELECT
	user_id,
	from_date,
	to_date,
	ST_X(from_location::geometry)::REAL as from_lon,
	ST_Y(from_location::geometry)::REAL as from_lat,
	en_route_order,
	tariff_per_km,
	from_address
FROM tow_driver_freely_available
WHERE user_id = $1;

-- name: GetFreelyAvailableLocations :many
SELECT
	id,
	tow_driver,
	ST_X(location::geometry)::REAL as lon,
	ST_Y(location::geometry)::REAL as lat,
	address
FROM tow_driver_freely_available_to_location_list
WHERE tow_driver = $1
ORDER BY id;

-- name: GetFreelyAvailableDrivers :many
SELECT
	tfa.user_id,
	tfa.from_date,
	tfa.to_date,
	ST_X(tfa.from_location::geometry)::REAL as from_lon,
	ST_Y(tfa.from_location::geometry)::REAL as from_lat,
	tfa.en_route_order,
	tfa.tariff_per_km,
	tfa.from_address,
	d.name,
	d.rating,
	d.profile_image,
	st_distance(
		tfa.from_location,
		st_setsrid(st_makepoint(@user_lon::REAL, @user_lat::REAL), 4326)::geometry
	)::real AS distance
FROM tow_driver_freely_available tfa
JOIN driver d ON tfa.user_id = d.user_id
WHERE
	tfa.from_date <= NOW() AND tfa.to_date >= NOW()
	AND (sqlc.narg('en_route_order')::BOOLEAN IS NULL OR tfa.en_route_order = sqlc.narg('en_route_order'))
	AND (sqlc.narg('min_tariff')::REAL IS NULL OR tfa.tariff_per_km >= sqlc.narg('min_tariff'))
	AND (sqlc.narg('max_tariff')::REAL IS NULL OR tfa.tariff_per_km <= sqlc.narg('max_tariff'))
ORDER BY distance;
