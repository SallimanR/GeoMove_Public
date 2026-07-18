-- name: CreateOrder :one
INSERT INTO "order" (
	customer_id,
	from_location,
	to_location,
	from_address,
	to_address,
	total_distance_meters,
	how_many_wheels_blocked,
	price_rubles,
	status
) VALUES (
	$1,
	ST_SetSRID(ST_MakePoint(@from_lon::REAL, @from_lat::REAL), 4326),
	ST_SetSRID(ST_MakePoint(@to_lon::REAL, @to_lat::REAL), 4326),
	@from_address,
	@to_address,
	$2,
	$3,
	$4,
	'forming'
)
RETURNING id;

-- name: GetOrderByID :one
SELECT
	id,
	created_at,
	updated_at,
	customer_id,
	driver_id,
	ST_X(from_location::geometry)::REAL AS from_lon,
	ST_Y(from_location::geometry)::REAL AS from_lat,
	ST_X(to_location::geometry)::REAL AS to_lon,
	ST_Y(to_location::geometry)::REAL AS to_lat,
	from_address,
	to_address,
	total_distance_meters,
	how_many_wheels_blocked,
	price_rubles,
	status,
	accepted_at,
	picked_up_at,
	completed_at,
	cancelled_at,
	cancellation_reason
FROM "order"
WHERE id = $1;

-- name: ListOrdersByCustomer :many
SELECT
	id,
	created_at,
	updated_at,
	customer_id,
	driver_id,
	ST_X(from_location::geometry)::REAL AS from_lon,
	ST_Y(from_location::geometry)::REAL AS from_lat,
	ST_X(to_location::geometry)::REAL AS to_lon,
	ST_Y(to_location::geometry)::REAL AS to_lat,
	from_address,
	to_address,
	total_distance_meters,
	how_many_wheels_blocked,
	price_rubles,
	status,
	accepted_at,
	picked_up_at,
	completed_at,
	cancelled_at,
	cancellation_reason
FROM "order"
WHERE customer_id = $1
ORDER BY created_at DESC;

-- name: ListOrdersByDriver :many
SELECT
	id,
	created_at,
	updated_at,
	customer_id,
	driver_id,
	ST_X(from_location::geometry)::REAL AS from_lon,
	ST_Y(from_location::geometry)::REAL AS from_lat,
	ST_X(to_location::geometry)::REAL AS to_lon,
	ST_Y(to_location::geometry)::REAL AS to_lat,
	from_address,
	to_address,
	total_distance_meters,
	how_many_wheels_blocked,
	price_rubles,
	status,
	accepted_at,
	picked_up_at,
	completed_at,
	cancelled_at,
	cancellation_reason
FROM "order"
WHERE driver_id = $1
ORDER BY created_at DESC;

-- name: UpdateOrderStatus :exec
UPDATE "order"
SET
	status = $2,
	accepted_at = COALESCE(sqlc.narg('accepted_at')::TIMESTAMP, accepted_at),
	picked_up_at = COALESCE(sqlc.narg('picked_up_at')::TIMESTAMP, picked_up_at),
	completed_at = COALESCE(sqlc.narg('completed_at')::TIMESTAMP, completed_at),
	cancelled_at = COALESCE(sqlc.narg('cancelled_at')::TIMESTAMP, cancelled_at),
	cancellation_reason = COALESCE(sqlc.narg('cancellation_reason'), cancellation_reason)
WHERE id = $1;

-- name: SetOrderDriver :exec
UPDATE "order"
SET driver_id = $2, status = 'accepted'
WHERE id = $1 AND status = 'pending';

-- name: UpdateOrderDetails :one
UPDATE "order"
SET
	from_location = ST_SetSRID(ST_MakePoint(@from_lon::REAL, @from_lat::REAL), 4326),
	to_location = ST_SetSRID(ST_MakePoint(@to_lon::REAL, @to_lat::REAL), 4326),
	from_address = @from_address,
	to_address = @to_address,
	total_distance_meters = $2,
	how_many_wheels_blocked = $3,
	price_rubles = $4,
	updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING
	id,
	created_at,
	updated_at,
	customer_id,
	driver_id,
	ST_X(from_location::geometry)::REAL AS from_lon,
	ST_Y(from_location::geometry)::REAL AS from_lat,
	ST_X(to_location::geometry)::REAL AS to_lon,
	ST_Y(to_location::geometry)::REAL AS to_lat,
	from_address,
	to_address,
	total_distance_meters,
	how_many_wheels_blocked,
	price_rubles,
	status,
	accepted_at,
	picked_up_at,
	completed_at,
	cancelled_at,
	cancellation_reason;
