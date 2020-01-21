--
-- PostgreSQL database dump
--

-- Dumped from database version 11.5 (Debian 11.5-1.pgdg90+1)
-- Dumped by pg_dump version 11.5 (Debian 11.5-1.pgdg100+1)

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

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: admins; Type: TABLE; Schema: public; Owner: USER
--

CREATE TABLE public.admins (
    id uuid NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.admins OWNER TO "USER";

--
-- Name: issues; Type: TABLE; Schema: public; Owner: USER
--

CREATE TABLE public.issues (
    id uuid NOT NULL,
    title text,
    experience_needed character varying(255) DEFAULT 'moderate'::character varying,
    expected_time character varying(255),
    language character varying(255),
    tech_stack character varying(255),
    github_id integer NOT NULL,
    number integer NOT NULL,
    labels character varying[],
    url text NOT NULL,
    body text,
    type character varying(255),
    repository_id uuid NOT NULL,
    project_id uuid NOT NULL,
    closed boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.issues OWNER TO "USER";

--
-- Name: projects; Type: TABLE; Schema: public; Owner: USER
--

CREATE TABLE public.projects (
    id uuid NOT NULL,
    display_name character varying(150) NOT NULL,
    first_color character varying(14) DEFAULT '#FF614C'::character varying NOT NULL,
    second_color character varying(14),
    description text NOT NULL,
    logo character varying(255) NOT NULL,
    link character varying(255) NOT NULL,
    setup_duration character varying(100),
    issues_count integer NOT NULL,
    tags character varying[] NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.projects OWNER TO "USER";

--
-- Name: repositories; Type: TABLE; Schema: public; Owner: USER
--

CREATE TABLE public.repositories (
    id uuid NOT NULL,
    repository_url character varying(255) NOT NULL,
    project_id uuid NOT NULL,
    issue_count integer DEFAULT 0 NOT NULL,
    last_parsed timestamp without time zone DEFAULT '1999-01-08 00:00:00'::timestamp without time zone NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    tags character varying[]
);


ALTER TABLE public.repositories OWNER TO "USER";

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: USER
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO "USER";

--
-- Name: admins admins_pkey; Type: CONSTRAINT; Schema: public; Owner: USER
--

ALTER TABLE ONLY public.admins
    ADD CONSTRAINT admins_pkey PRIMARY KEY (id);


--
-- Name: issues issues_pkey; Type: CONSTRAINT; Schema: public; Owner: USER
--

ALTER TABLE ONLY public.issues
    ADD CONSTRAINT issues_pkey PRIMARY KEY (id);


--
-- Name: projects projects_pkey; Type: CONSTRAINT; Schema: public; Owner: USER
--

ALTER TABLE ONLY public.projects
    ADD CONSTRAINT projects_pkey PRIMARY KEY (id);


--
-- Name: repositories repositories_pkey; Type: CONSTRAINT; Schema: public; Owner: USER
--

ALTER TABLE ONLY public.repositories
    ADD CONSTRAINT repositories_pkey PRIMARY KEY (id);


--
-- Name: index_issue_experience_needed; Type: INDEX; Schema: public; Owner: USER
--

CREATE INDEX index_issue_experience_needed ON public.issues USING btree (experience_needed);


--
-- Name: index_issue_language; Type: INDEX; Schema: public; Owner: USER
--

CREATE INDEX index_issue_language ON public.issues USING btree (language);


--
-- Name: index_issue_type; Type: INDEX; Schema: public; Owner: USER
--

CREATE INDEX index_issue_type ON public.issues USING btree (type);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: USER
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: issues issues_projects_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: USER
--

ALTER TABLE ONLY public.issues
    ADD CONSTRAINT issues_projects_id_fk FOREIGN KEY (project_id) REFERENCES public.projects(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: issues issues_repositories_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: USER
--

ALTER TABLE ONLY public.issues
    ADD CONSTRAINT issues_repositories_id_fk FOREIGN KEY (repository_id) REFERENCES public.repositories(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: repositories repositories_projects_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: USER
--

ALTER TABLE ONLY public.repositories
    ADD CONSTRAINT repositories_projects_id_fk FOREIGN KEY (project_id) REFERENCES public.projects(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

