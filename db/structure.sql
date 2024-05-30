--
-- PostgreSQL database dump
--

-- Dumped from database version 15.7 (Homebrew)
-- Dumped by pg_dump version 15.7 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE ONLY public.posts DROP CONSTRAINT posts_feed_id_fkey;
ALTER TABLE ONLY public.feeds DROP CONSTRAINT feeds_user_id_fkey;
ALTER TABLE ONLY public.feed_follows DROP CONSTRAINT feed_follows_user_id_fkey;
ALTER TABLE ONLY public.feed_follows DROP CONSTRAINT feed_follows_feed_id_fkey;
ALTER TABLE ONLY public.users DROP CONSTRAINT users_pkey;
ALTER TABLE ONLY public.posts DROP CONSTRAINT posts_url_key;
ALTER TABLE ONLY public.posts DROP CONSTRAINT posts_pkey;
ALTER TABLE ONLY public.goose_db_version DROP CONSTRAINT goose_db_version_pkey;
ALTER TABLE ONLY public.feeds DROP CONSTRAINT feeds_url_key;
ALTER TABLE ONLY public.feeds DROP CONSTRAINT feeds_pkey;
ALTER TABLE ONLY public.feed_follows DROP CONSTRAINT feed_follows_pkey;
ALTER TABLE ONLY public.feed_follows DROP CONSTRAINT feed_follows_feed_id_user_id_key;
ALTER TABLE public.goose_db_version ALTER COLUMN id DROP DEFAULT;
DROP TABLE public.users;
DROP TABLE public.posts;
DROP SEQUENCE public.goose_db_version_id_seq;
DROP TABLE public.goose_db_version;
DROP TABLE public.feeds;
DROP TABLE public.feed_follows;
SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: feed_follows; Type: TABLE; Schema: public; Owner: app_scrapygo
--

CREATE TABLE public.feed_follows (
    id uuid NOT NULL,
    feed_id uuid NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


ALTER TABLE public.feed_follows OWNER TO app_scrapygo;

--
-- Name: feeds; Type: TABLE; Schema: public; Owner: app_scrapygo
--

CREATE TABLE public.feeds (
    id uuid NOT NULL,
    name character varying NOT NULL,
    url character varying NOT NULL,
    user_id uuid NOT NULL,
    last_fetched_at timestamp with time zone,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


ALTER TABLE public.feeds OWNER TO app_scrapygo;

--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: app_scrapygo
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);


ALTER TABLE public.goose_db_version OWNER TO app_scrapygo;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: app_scrapygo
--

CREATE SEQUENCE public.goose_db_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.goose_db_version_id_seq OWNER TO app_scrapygo;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: app_scrapygo
--

ALTER SEQUENCE public.goose_db_version_id_seq OWNED BY public.goose_db_version.id;


--
-- Name: posts; Type: TABLE; Schema: public; Owner: app_scrapygo
--

CREATE TABLE public.posts (
    id uuid NOT NULL,
    feed_id uuid NOT NULL,
    title character varying(500) NOT NULL,
    description text,
    url character varying NOT NULL,
    published_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


ALTER TABLE public.posts OWNER TO app_scrapygo;

--
-- Name: users; Type: TABLE; Schema: public; Owner: app_scrapygo
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    name character varying NOT NULL,
    api_key character varying DEFAULT encode(sha256(((random())::text)::bytea), 'hex'::text) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


ALTER TABLE public.users OWNER TO app_scrapygo;

--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: app_scrapygo
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Name: feed_follows feed_follows_feed_id_user_id_key; Type: CONSTRAINT; Schema: public; Owner: app_scrapygo
--

ALTER TABLE ONLY public.feed_follows
    ADD CONSTRAINT feed_follows_feed_id_user_id_key UNIQUE (feed_id, user_id);


--
-- Name: feed_follows feed_follows_pkey; Type: CONSTRAINT; Schema: public; Owner: app_scrapygo
--

ALTER TABLE ONLY public.feed_follows
    ADD CONSTRAINT feed_follows_pkey PRIMARY KEY (id);


--
-- Name: feeds feeds_pkey; Type: CONSTRAINT; Schema: public; Owner: app_scrapygo
--

ALTER TABLE ONLY public.feeds
    ADD CONSTRAINT feeds_pkey PRIMARY KEY (id);


--
-- Name: feeds feeds_url_key; Type: CONSTRAINT; Schema: public; Owner: app_scrapygo
--

ALTER TABLE ONLY public.feeds
    ADD CONSTRAINT feeds_url_key UNIQUE (url);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: app_scrapygo
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: app_scrapygo
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- Name: posts posts_url_key; Type: CONSTRAINT; Schema: public; Owner: app_scrapygo
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_url_key UNIQUE (url);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: app_scrapygo
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: feed_follows feed_follows_feed_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: app_scrapygo
--

ALTER TABLE ONLY public.feed_follows
    ADD CONSTRAINT feed_follows_feed_id_fkey FOREIGN KEY (feed_id) REFERENCES public.feeds(id) ON DELETE CASCADE;


--
-- Name: feed_follows feed_follows_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: app_scrapygo
--

ALTER TABLE ONLY public.feed_follows
    ADD CONSTRAINT feed_follows_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: feeds feeds_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: app_scrapygo
--

ALTER TABLE ONLY public.feeds
    ADD CONSTRAINT feeds_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: posts posts_feed_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: app_scrapygo
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_feed_id_fkey FOREIGN KEY (feed_id) REFERENCES public.feeds(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

