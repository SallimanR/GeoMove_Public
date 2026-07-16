\restrict dZ882t0xCG2xGijeFniuK5Rml0oZiK9RgUbN8UjPfFrMM4QeDiRWut3veA0aHif

-- Dumped from database version 18.4 (Debian 18.4-1.pgdg13+1)
-- Dumped by pg_dump version 18.4

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
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: driver; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.driver (
    user_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    name text NOT NULL,
    profile_image text,
    work_starts time without time zone,
    work_ends time without time zone,
    is_available boolean DEFAULT true NOT NULL,
    last_seen timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    location public.geography(Point,4326) NOT NULL,
    rating real
);


--
-- Name: moving_driver; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.moving_driver (
    driver_id bigint NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    realtime_location public.geography(Point,4326),
    travel_time time without time zone NOT NULL,
    path_meters integer NOT NULL,
    coarse_h3 public.h3index GENERATED ALWAYS AS (public.h3_latlng_to_cell(realtime_location, 2)) STORED,
    destination_location public.geography(Point,4326),
    destination_time timestamp without time zone
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
    token_hash text NOT NULL,
    session_id uuid DEFAULT public.uuid_generate_v4(),
    user_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    roles text[] DEFAULT '{}'::text[] NOT NULL
);


--
-- Name: tow_driver_freely_available; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tow_driver_freely_available (
    user_id bigint NOT NULL,
    from_date timestamp without time zone NOT NULL,
    to_date timestamp without time zone NOT NULL,
    from_location public.geography(Point,4326) NOT NULL,
    from_address text DEFAULT ''::text NOT NULL,
    en_route_order boolean,
    tariff_per_km real
);


--
-- Name: tow_driver_freely_available_to_location_list; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tow_driver_freely_available_to_location_list (
    id bigint NOT NULL,
    tow_driver bigint CONSTRAINT tow_driver_freely_available_to_location_lis_tow_driver_not_null NOT NULL,
    location public.geography(Point,4326) NOT NULL,
    address text DEFAULT ''::text NOT NULL
);


--
-- Name: tow_driver_freely_available_to_location_list_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tow_driver_freely_available_to_location_list_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tow_driver_freely_available_to_location_list_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tow_driver_freely_available_to_location_list_id_seq OWNED BY public.tow_driver_freely_available_to_location_list.id;


--
-- Name: user; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."user" (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    phone text,
    email text,
    profile_image text,
    roles text[] DEFAULT '{user}'::text[] NOT NULL
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
-- Name: user_oauth_links; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_oauth_links (
    user_id bigint NOT NULL,
    provider text NOT NULL,
    provider_id text NOT NULL
);


--
-- Name: tow_driver_freely_available_to_location_list id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tow_driver_freely_available_to_location_list ALTER COLUMN id SET DEFAULT nextval('public.tow_driver_freely_available_to_location_list_id_seq'::regclass);


--
-- Name: moving_driver moving_driver_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.moving_driver
    ADD CONSTRAINT moving_driver_pkey PRIMARY KEY (driver_id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: session session_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_pkey PRIMARY KEY (token_hash);


--
-- Name: tow_driver_freely_available tow_driver_freely_available_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tow_driver_freely_available
    ADD CONSTRAINT tow_driver_freely_available_pkey PRIMARY KEY (user_id);


--
-- Name: tow_driver_freely_available_to_location_list tow_driver_freely_available_to_location_list_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tow_driver_freely_available_to_location_list
    ADD CONSTRAINT tow_driver_freely_available_to_location_list_pkey PRIMARY KEY (id);


--
-- Name: user user_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_email_key UNIQUE (email);


--
-- Name: user_oauth_links user_oauth_links_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_oauth_links
    ADD CONSTRAINT user_oauth_links_pkey PRIMARY KEY (user_id, provider);


--
-- Name: user user_phone_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_phone_key UNIQUE (phone);


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- Name: user user_profile_image_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_profile_image_key UNIQUE (profile_image);


--
-- Name: idx_destination_location; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_destination_location ON public.moving_driver USING gist (destination_location);


--
-- Name: idx_destination_time; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_destination_time ON public.moving_driver USING btree (destination_time);


--
-- Name: idx_driver_location_geom; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_driver_location_geom ON public.driver USING gist (location);


--
-- Name: idx_driver_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_driver_user_id ON public.driver USING btree (user_id);


--
-- Name: idx_moving_driver_coarse_h3; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_moving_driver_coarse_h3 ON public.moving_driver USING btree (coarse_h3);


--
-- Name: idx_moving_driver_location; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_moving_driver_location ON public.moving_driver USING gist (realtime_location);


--
-- Name: idx_session_expires_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_session_expires_at ON public.session USING btree (expires_at);


--
-- Name: idx_session_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_session_user_id ON public.session USING btree (user_id);


--
-- Name: idx_tow_driv_fa_loc_geom; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tow_driv_fa_loc_geom ON public.tow_driver_freely_available USING gist (from_location);


--
-- Name: idx_tow_drv_fa_loc_list_fk; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tow_drv_fa_loc_list_fk ON public.tow_driver_freely_available_to_location_list USING btree (tow_driver);


--
-- Name: idx_tow_drv_fa_loc_list_geom; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tow_drv_fa_loc_list_geom ON public.tow_driver_freely_available_to_location_list USING gist (location);


--
-- Name: idx_tow_fa_drv_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_tow_fa_drv_user_id ON public.tow_driver_freely_available USING btree (user_id);


--
-- Name: idx_user_oauth_links_provider; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_oauth_links_provider ON public.user_oauth_links USING btree (provider, provider_id);


--
-- Name: driver driver_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver
    ADD CONSTRAINT driver_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: moving_driver moving_driver_driver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.moving_driver
    ADD CONSTRAINT moving_driver_driver_id_fkey FOREIGN KEY (driver_id) REFERENCES public.driver(user_id) ON DELETE CASCADE;


--
-- Name: session session_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: tow_driver_freely_available_to_location_list tow_driver_freely_available_to_location_list_tow_driver_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tow_driver_freely_available_to_location_list
    ADD CONSTRAINT tow_driver_freely_available_to_location_list_tow_driver_fkey FOREIGN KEY (tow_driver) REFERENCES public.tow_driver_freely_available(user_id) ON DELETE CASCADE;


--
-- Name: tow_driver_freely_available tow_driver_freely_available_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tow_driver_freely_available
    ADD CONSTRAINT tow_driver_freely_available_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.driver(user_id) ON DELETE CASCADE;


--
-- Name: user_oauth_links user_oauth_links_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_oauth_links
    ADD CONSTRAINT user_oauth_links_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict dZ882t0xCG2xGijeFniuK5Rml0oZiK9RgUbN8UjPfFrMM4QeDiRWut3veA0aHif


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20260110181007'),
    ('20260110181146'),
    ('20260322214503'),
    ('20260323100924'),
    ('20260705102447'),
    ('20260706094702'),
    ('20260710043141'),
    ('20260713162812'),
    ('20260713193617');
