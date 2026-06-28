CREATE OR REPLACE FUNCTION filter_driver_points_optimized(
    z integer, 
    x integer, 
    y integer, 
    query_params json
) RETURNS bytea AS $$
DECLARE
    mvt bytea;
    tile_envelope_3857 geometry;
    param_work_starts timetz;
    param_work_ends timetz;
    param_rating int;
    sql_text text;
BEGIN
    -- 1. Extract parameters ONCE
    param_work_starts := (query_params->>'work_starts')::timetz;
    param_work_ends := (query_params->>'work_ends')::timetz;
    param_rating := (query_params->>'rating')::int;

    -- 2. Build a dynamic WHERE clause
    sql_text := '
    SELECT ST_AsMVT(tile, ''filter_driver_points'', 4096, ''geom'')
    FROM (
        SELECT
            ST_AsMVTGeom(
                ST_Transform(ST_CurveToLine(d.location::geometry), 3857),
                ST_TileEnvelope($1, $2, $3), -- Use tile envelope directly
                4096, 64, true
            ) AS geom,
            d.city_id,
            d.work_starts,
            d.work_ends,
            d.rating
        FROM public.drivers d
        WHERE d.location && ST_Transform(ST_TileEnvelope($1, $2, $3), 4326)';

    -- 3. Conditionally add filters
    IF param_work_starts IS NOT NULL THEN
        sql_text := sql_text || ' AND d.work_starts <= $4';
    END IF;
    
    IF param_work_ends IS NOT NULL THEN
        sql_text := sql_text || ' AND d.work_ends <= $5';
    END IF;
    
    IF param_rating IS NOT NULL THEN
        sql_text := sql_text || ' AND d.rating >= $6';
    END IF;

    sql_text := sql_text || ') AS tile WHERE geom IS NOT NULL';

    -- 4. Execute with only necessary parameters
    IF param_work_starts IS NULL AND param_work_ends IS NULL AND param_rating IS NULL THEN
        EXECUTE sql_text INTO mvt USING z, x, y;
    ELSIF param_work_starts IS NOT NULL AND param_work_ends IS NULL AND param_rating IS NULL THEN
        EXECUTE sql_text INTO mvt USING z, x, y, param_work_starts;
    -- ... handle other combinations (could use a CASE or more IF statements)
    -- For simplicity in this example, let's assume all params might be provided:
    ELSE
        EXECUTE sql_text INTO mvt USING z, x, y, param_work_starts, param_work_ends, param_rating;
    END IF;

    RETURN mvt;
END;
$$ LANGUAGE plpgsql IMMUTABLE STRICT PARALLEL SAFE;
