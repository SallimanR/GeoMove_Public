\restrict k8kf7C2qj4JDMBZZ6kj2hFnIQljpJ6OPw10cakUZfC2rJjiOvETdL2zwgeLVxsl

-- Dumped from database version 18.1 (Debian 18.1-1.pgdg13+2)
-- Dumped by pg_dump version 18.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: h3; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS h3 WITH SCHEMA public;


--
-- Name: EXTENSION h3; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION h3 IS 'H3 bindings for PostgreSQL';


--
-- Name: postgis; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA public;


--
-- Name: EXTENSION postgis; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION postgis IS 'PostGIS geometry and geography spatial types and functions';


--
-- Name: postgis_raster; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis_raster WITH SCHEMA public;


--
-- Name: EXTENSION postgis_raster; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION postgis_raster IS 'PostGIS raster types and functions';


--
-- Name: h3_postgis; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS h3_postgis WITH SCHEMA public;


--
-- Name: EXTENSION h3_postgis; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION h3_postgis IS 'H3 PostGIS integration';


--
-- Name: update_driver_location_h3_indexes(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_driver_location_h3_indexes() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.h3_res5 = h3_latlng_to_cell(NEW.location::GEOMETRY, 5);
    NEW.h3_res6 = h3_latlng_to_cell(NEW.location::GEOMETRY, 6);
    NEW.h3_res7 = h3_latlng_to_cell(NEW.location::GEOMETRY, 7);

	NEW.updated_at = CURRENT_TIMESTAMP;

	RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: access_token; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.access_token (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    session_id bigint NOT NULL,
    token_hash text NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    revoked_at timestamp without time zone
);


--
-- Name: driver; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.driver (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    work_starts time without time zone,
    work_ends time without time zone,
    rating real,
    is_available boolean DEFAULT true,
    last_seen timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    location public.geography(Point,4326) NOT NULL,
    h3_res5 public.h3index,
    h3_res6 public.h3index,
    h3_res7 public.h3index
);


--
-- Name: driver_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.driver ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.driver_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: driver_realtime; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.driver_realtime (
    driver_id bigint NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    realtime_location public.geography(Point,4326),
    average_speed real,
    predicted_bearing real,
    coarse_h3 public.h3index GENERATED ALWAYS AS (public.h3_latlng_to_cell(realtime_location, 2)) STORED,
    destination_location public.geography(Point,4326),
    destination_time timestamp without time zone
);


--
-- Name: role; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.role (
    id integer NOT NULL,
    name text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: session; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.session (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone
);


--
-- Name: user; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."user" (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone
);


--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public."user" ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: user_role; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_role (
    user_id bigint NOT NULL,
    role_id integer NOT NULL,
    assigned_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: access_token access_token_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.access_token
    ADD CONSTRAINT access_token_pkey PRIMARY KEY (id);


--
-- Name: access_token access_token_token_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.access_token
    ADD CONSTRAINT access_token_token_hash_key UNIQUE (token_hash);


--
-- Name: driver driver_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver
    ADD CONSTRAINT driver_pkey PRIMARY KEY (id);


--
-- Name: driver_realtime driver_realtime_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver_realtime
    ADD CONSTRAINT driver_realtime_pkey PRIMARY KEY (driver_id);


--
-- Name: driver driver_user_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver
    ADD CONSTRAINT driver_user_id_key UNIQUE (user_id);


--
-- Name: role role_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role
    ADD CONSTRAINT role_name_key UNIQUE (name);


--
-- Name: role role_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role
    ADD CONSTRAINT role_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: session session_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_pkey PRIMARY KEY (id);


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- Name: user_role user_role_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_role
    ADD CONSTRAINT user_role_pkey PRIMARY KEY (user_id, role_id);


--
-- Name: idx_destination_location; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_destination_location ON public.driver_realtime USING gist (destination_location);


--
-- Name: idx_destination_time; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_destination_time ON public.driver_realtime USING btree (destination_time);


--
-- Name: idx_driver_h3_res5; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_driver_h3_res5 ON public.driver USING btree (h3_res5);


--
-- Name: idx_driver_h3_res6; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_driver_h3_res6 ON public.driver USING btree (h3_res6);


--
-- Name: idx_driver_h3_res7; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_driver_h3_res7 ON public.driver USING btree (h3_res7);


--
-- Name: idx_driver_location_geom; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_driver_location_geom ON public.driver USING gist (location);


--
-- Name: idx_driver_realtime_coarse_h3; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_driver_realtime_coarse_h3 ON public.driver_realtime USING btree (coarse_h3);


--
-- Name: idx_driver_realtime_location; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_driver_realtime_location ON public.driver_realtime USING gist (realtime_location);


--
-- Name: idx_driver_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_driver_user_id ON public.driver USING btree (user_id);


--
-- Name: idx_token_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_token_hash ON public.access_token USING btree (token_hash);


--
-- Name: idx_user_roles_role; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_role ON public.user_role USING btree (role_id);


--
-- Name: idx_user_roles_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_user ON public.user_role USING btree (user_id);


--
-- Name: driver trigger_driver_location_update; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_driver_location_update BEFORE INSERT OR UPDATE OF location ON public.driver FOR EACH ROW EXECUTE FUNCTION public.update_driver_location_h3_indexes();


--
-- Name: access_token access_token_session_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.access_token
    ADD CONSTRAINT access_token_session_id_fkey FOREIGN KEY (session_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: driver_realtime driver_realtime_driver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver_realtime
    ADD CONSTRAINT driver_realtime_driver_id_fkey FOREIGN KEY (driver_id) REFERENCES public.driver(id) ON DELETE CASCADE;


--
-- Name: driver driver_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver
    ADD CONSTRAINT driver_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: session session_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: user_role user_role_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_role
    ADD CONSTRAINT user_role_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.role(id) ON DELETE CASCADE;


--
-- Name: user_role user_role_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_role
    ADD CONSTRAINT user_role_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict k8kf7C2qj4JDMBZZ6kj2hFnIQljpJ6OPw10cakUZfC2rJjiOvETdL2zwgeLVxsl


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20260110181007'),
    ('20260110181146'),
    ('20260322214503'),
    ('20260323095635'),
    ('20260323100924'),
    ('20260323101002'),
    ('20260323101003');
