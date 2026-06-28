-- migrate:up
CREATE OR REPLACE FUNCTION generate_geography()
RETURNS TRIGGER AS $$
DECLARE
    _city_id int;
    _state_id int;
    _city_data city;
    _state_data state;
BEGIN
    -- Skip if location hasn't changed (unless null references)
    IF TG_OP = 'UPDATE' AND OLD.location IS NOT DISTINCT FROM NEW.location 
       AND OLD.city_id IS NOT NULL AND OLD.state_id IS NOT NULL THEN
        RETURN NEW;
    END IF;
    
    -- Compute H3 indexes
    NEW.h3_res8 := h3_geo_to_h3(NEW.location, 8);
    NEW.h3_res9 := h3_geo_to_h3(NEW.location, 9);
    NEW.h3_res10 := h3_geo_to_h3(NEW.location, 10);
    
    -- Find state using H3 first (much faster than ST_Within)
    SELECT s.id, s.* INTO _state_data
    FROM state s
    WHERE NEW.h3_res8 = ANY(s.h3_res8)
      AND ST_Within(NEW.location, s.geometry)
    LIMIT 1;
    
    -- Find city
    SELECT c.id, c.* INTO _city_data
    FROM city c
    WHERE NEW.h3_res9 = ANY(c.h3_res9)
      AND ST_Within(NEW.location, c.geometry)
    LIMIT 1;
    
    -- Set references
    NEW.state_id := _state_data.id;
    NEW.city_id := _city_data.id;
    
    -- Cache frequently accessed data
    IF _city_data.id IS NOT NULL THEN
        NEW.city_cache := jsonb_build_object(
            'name', _city_data.name,
            'timezone', _city_data.timezone,
            'population', _city_data.population
        );
    END IF;
    
    IF _state_data.id IS NOT NULL THEN
        NEW.state_cache := jsonb_build_object(
            'name', _state_data.name,
            'iso_code', _state_data.iso_code
        );
    END IF;
    
    RETURN NEW;
EXCEPTION
    WHEN OTHERS THEN
        RAISE LOG 'Failed to generate geography for driver %: %', COALESCE(NEW.id, '1'), SQLERRM;
        RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- migrate:down
DROP FUNCTION generate_geography
