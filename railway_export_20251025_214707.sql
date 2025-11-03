--
-- PostgreSQL database dump
--

\restrict 7LGbFidJ9dXINV6qGK8WtxaikiEKfaIeemyhRVZIRtRconJiZjgAIhHVKud6sCe

-- Dumped from database version 15.14
-- Dumped by pg_dump version 15.14

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

ALTER TABLE IF EXISTS ONLY public.user_consents DROP CONSTRAINT IF EXISTS user_consents_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.relationship_analyses DROP CONSTRAINT IF EXISTS relationship_analyses_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.relationship_analyses DROP CONSTRAINT IF EXISTS relationship_analyses_person_id_fkey;
ALTER TABLE IF EXISTS ONLY public.refresh_tokens DROP CONSTRAINT IF EXISTS refresh_tokens_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.reflections DROP CONSTRAINT IF EXISTS reflections_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.push_tokens DROP CONSTRAINT IF EXISTS push_tokens_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.people DROP CONSTRAINT IF EXISTS people_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.people DROP CONSTRAINT IF EXISTS people_relationship_status_id_fkey;
ALTER TABLE IF EXISTS ONLY public.people DROP CONSTRAINT IF EXISTS people_intention_id_fkey;
ALTER TABLE IF EXISTS ONLY public.people DROP CONSTRAINT IF EXISTS people_energy_pattern_id_fkey;
ALTER TABLE IF EXISTS ONLY public.people DROP CONSTRAINT IF EXISTS people_communication_method_id_fkey;
ALTER TABLE IF EXISTS ONLY public.people DROP CONSTRAINT IF EXISTS people_category_id_fkey;
ALTER TABLE IF EXISTS ONLY public.nudges DROP CONSTRAINT IF EXISTS nudges_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.nudges DROP CONSTRAINT IF EXISTS nudges_person_id_fkey;
ALTER TABLE IF EXISTS ONLY public.nudges DROP CONSTRAINT IF EXISTS nudges_analysis_id_fkey;
ALTER TABLE IF EXISTS ONLY public.interactions DROP CONSTRAINT IF EXISTS interactions_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.interactions DROP CONSTRAINT IF EXISTS interactions_person_id_fkey;
ALTER TABLE IF EXISTS ONLY public.events DROP CONSTRAINT IF EXISTS events_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.data_exports DROP CONSTRAINT IF EXISTS data_exports_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.daily_metrics DROP CONSTRAINT IF EXISTS daily_metrics_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.categories DROP CONSTRAINT IF EXISTS categories_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.auth_providers DROP CONSTRAINT IF EXISTS auth_providers_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.audit_logs DROP CONSTRAINT IF EXISTS audit_logs_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.ai_analysis_jobs DROP CONSTRAINT IF EXISTS ai_analysis_jobs_user_id_fkey;
DROP TRIGGER IF EXISTS update_users_updated_at ON public.users;
DROP TRIGGER IF EXISTS update_user_consents_updated_at ON public.user_consents;
DROP TRIGGER IF EXISTS update_refresh_tokens_updated_at ON public.refresh_tokens;
DROP TRIGGER IF EXISTS update_reflections_updated_at ON public.reflections;
DROP TRIGGER IF EXISTS update_push_tokens_updated_at ON public.push_tokens;
DROP TRIGGER IF EXISTS update_person_last_interaction_trigger ON public.interactions;
DROP TRIGGER IF EXISTS update_people_updated_at ON public.people;
DROP TRIGGER IF EXISTS update_nudges_updated_at ON public.nudges;
DROP TRIGGER IF EXISTS update_interactions_updated_at ON public.interactions;
DROP TRIGGER IF EXISTS update_events_updated_at ON public.events;
DROP TRIGGER IF EXISTS update_data_exports_updated_at ON public.data_exports;
DROP TRIGGER IF EXISTS update_daily_metrics_updated_at ON public.daily_metrics;
DROP TRIGGER IF EXISTS update_auth_providers_updated_at ON public.auth_providers;
DROP TRIGGER IF EXISTS update_audit_logs_updated_at ON public.audit_logs;
DROP TRIGGER IF EXISTS update_ai_analysis_jobs_updated_at ON public.ai_analysis_jobs;
DROP INDEX IF EXISTS public.idx_users_username;
DROP INDEX IF EXISTS public.idx_users_email;
DROP INDEX IF EXISTS public.idx_users_deleted_at;
DROP INDEX IF EXISTS public.idx_users_data_residency;
DROP INDEX IF EXISTS public.idx_user_consents_user_id;
DROP INDEX IF EXISTS public.idx_user_consents_type;
DROP INDEX IF EXISTS public.idx_user_consents_granted;
DROP INDEX IF EXISTS public.idx_relationship_analyses_user_person;
DROP INDEX IF EXISTS public.idx_relationship_analyses_user_id;
DROP INDEX IF EXISTS public.idx_relationship_analyses_person_id;
DROP INDEX IF EXISTS public.idx_relationship_analyses_deleted_at;
DROP INDEX IF EXISTS public.idx_relationship_analyses_analyzed_at;
DROP INDEX IF EXISTS public.idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS public.idx_refresh_tokens_token;
DROP INDEX IF EXISTS public.idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS public.idx_reflections_user_id;
DROP INDEX IF EXISTS public.idx_reflections_user_date;
DROP INDEX IF EXISTS public.idx_reflections_mood;
DROP INDEX IF EXISTS public.idx_reflections_completed_at;
DROP INDEX IF EXISTS public.idx_push_tokens_user_id;
DROP INDEX IF EXISTS public.idx_push_tokens_token;
DROP INDEX IF EXISTS public.idx_push_tokens_active;
DROP INDEX IF EXISTS public.idx_people_user_id;
DROP INDEX IF EXISTS public.idx_people_relationship_status_id;
DROP INDEX IF EXISTS public.idx_people_next_reminder;
DROP INDEX IF EXISTS public.idx_people_last_interaction;
DROP INDEX IF EXISTS public.idx_people_health_score;
DROP INDEX IF EXISTS public.idx_people_energy_pattern_id;
DROP INDEX IF EXISTS public.idx_people_comm_method_id;
DROP INDEX IF EXISTS public.idx_people_category_id;
DROP INDEX IF EXISTS public.idx_nudges_user_id;
DROP INDEX IF EXISTS public.idx_nudges_type;
DROP INDEX IF EXISTS public.idx_nudges_timing;
DROP INDEX IF EXISTS public.idx_nudges_status;
DROP INDEX IF EXISTS public.idx_nudges_source;
DROP INDEX IF EXISTS public.idx_nudges_seen;
DROP INDEX IF EXISTS public.idx_nudges_priority;
DROP INDEX IF EXISTS public.idx_nudges_person_id;
DROP INDEX IF EXISTS public.idx_nudges_expires_at;
DROP INDEX IF EXISTS public.idx_nudges_analysis_id;
DROP INDEX IF EXISTS public.idx_nudges_acted_on;
DROP INDEX IF EXISTS public.idx_interactions_user_id;
DROP INDEX IF EXISTS public.idx_interactions_user_date;
DROP INDEX IF EXISTS public.idx_interactions_person_id;
DROP INDEX IF EXISTS public.idx_interactions_interaction_at;
DROP INDEX IF EXISTS public.idx_interactions_energy_impact;
DROP INDEX IF EXISTS public.idx_events_user_type_date;
DROP INDEX IF EXISTS public.idx_events_user_id;
DROP INDEX IF EXISTS public.idx_events_session_id;
DROP INDEX IF EXISTS public.idx_events_event_type;
DROP INDEX IF EXISTS public.idx_events_created_at;
DROP INDEX IF EXISTS public.idx_data_exports_user_id;
DROP INDEX IF EXISTS public.idx_data_exports_status;
DROP INDEX IF EXISTS public.idx_data_exports_expires_at;
DROP INDEX IF EXISTS public.idx_daily_metrics_user_id;
DROP INDEX IF EXISTS public.idx_daily_metrics_user_date;
DROP INDEX IF EXISTS public.idx_daily_metrics_date;
DROP INDEX IF EXISTS public.idx_auth_providers_user_id;
DROP INDEX IF EXISTS public.idx_auth_providers_provider;
DROP INDEX IF EXISTS public.idx_audit_logs_user_id;
DROP INDEX IF EXISTS public.idx_audit_logs_request_id;
DROP INDEX IF EXISTS public.idx_audit_logs_entity;
DROP INDEX IF EXISTS public.idx_audit_logs_created_at;
DROP INDEX IF EXISTS public.idx_audit_logs_action;
DROP INDEX IF EXISTS public.idx_ai_analysis_jobs_user_id;
DROP INDEX IF EXISTS public.idx_ai_analysis_jobs_status;
DROP INDEX IF EXISTS public.idx_ai_analysis_jobs_priority;
DROP INDEX IF EXISTS public.idx_ai_analysis_jobs_deleted_at;
DROP INDEX IF EXISTS public.idx_ai_analysis_jobs_created_at;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_username_key;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE IF EXISTS ONLY public.user_consents DROP CONSTRAINT IF EXISTS user_consents_pkey;
ALTER TABLE IF EXISTS ONLY public.categories DROP CONSTRAINT IF EXISTS uq_categories_user_name;
ALTER TABLE IF EXISTS ONLY public.relationship_statuses DROP CONSTRAINT IF EXISTS relationship_statuses_pkey;
ALTER TABLE IF EXISTS ONLY public.relationship_statuses DROP CONSTRAINT IF EXISTS relationship_statuses_name_key;
ALTER TABLE IF EXISTS ONLY public.relationship_analyses DROP CONSTRAINT IF EXISTS relationship_analyses_pkey;
ALTER TABLE IF EXISTS ONLY public.refresh_tokens DROP CONSTRAINT IF EXISTS refresh_tokens_token_key;
ALTER TABLE IF EXISTS ONLY public.refresh_tokens DROP CONSTRAINT IF EXISTS refresh_tokens_pkey;
ALTER TABLE IF EXISTS ONLY public.reflections DROP CONSTRAINT IF EXISTS reflections_pkey;
ALTER TABLE IF EXISTS ONLY public.push_tokens DROP CONSTRAINT IF EXISTS push_tokens_token_key;
ALTER TABLE IF EXISTS ONLY public.push_tokens DROP CONSTRAINT IF EXISTS push_tokens_pkey;
ALTER TABLE IF EXISTS ONLY public.people DROP CONSTRAINT IF EXISTS people_pkey;
ALTER TABLE IF EXISTS ONLY public.nudges DROP CONSTRAINT IF EXISTS nudges_pkey;
ALTER TABLE IF EXISTS ONLY public.interactions DROP CONSTRAINT IF EXISTS interactions_pkey;
ALTER TABLE IF EXISTS ONLY public.intentions DROP CONSTRAINT IF EXISTS intentions_pkey;
ALTER TABLE IF EXISTS ONLY public.intentions DROP CONSTRAINT IF EXISTS intentions_name_key;
ALTER TABLE IF EXISTS ONLY public.events DROP CONSTRAINT IF EXISTS events_pkey;
ALTER TABLE IF EXISTS ONLY public.energy_patterns DROP CONSTRAINT IF EXISTS energy_patterns_pkey;
ALTER TABLE IF EXISTS ONLY public.energy_patterns DROP CONSTRAINT IF EXISTS energy_patterns_name_key;
ALTER TABLE IF EXISTS ONLY public.data_exports DROP CONSTRAINT IF EXISTS data_exports_pkey;
ALTER TABLE IF EXISTS ONLY public.daily_metrics DROP CONSTRAINT IF EXISTS daily_metrics_pkey;
ALTER TABLE IF EXISTS ONLY public.communication_methods DROP CONSTRAINT IF EXISTS communication_methods_pkey;
ALTER TABLE IF EXISTS ONLY public.communication_methods DROP CONSTRAINT IF EXISTS communication_methods_name_key;
ALTER TABLE IF EXISTS ONLY public.categories DROP CONSTRAINT IF EXISTS categories_pkey;
ALTER TABLE IF EXISTS ONLY public.auth_providers DROP CONSTRAINT IF EXISTS auth_providers_provider_provider_id_key;
ALTER TABLE IF EXISTS ONLY public.auth_providers DROP CONSTRAINT IF EXISTS auth_providers_pkey;
ALTER TABLE IF EXISTS ONLY public.audit_logs DROP CONSTRAINT IF EXISTS audit_logs_pkey;
ALTER TABLE IF EXISTS ONLY public.ai_analysis_jobs DROP CONSTRAINT IF EXISTS ai_analysis_jobs_pkey;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.user_consents;
DROP TABLE IF EXISTS public.relationship_statuses;
DROP TABLE IF EXISTS public.relationship_analyses;
DROP TABLE IF EXISTS public.refresh_tokens;
DROP TABLE IF EXISTS public.reflections;
DROP TABLE IF EXISTS public.push_tokens;
DROP TABLE IF EXISTS public.people;
DROP TABLE IF EXISTS public.nudges;
DROP TABLE IF EXISTS public.interactions;
DROP TABLE IF EXISTS public.intentions;
DROP TABLE IF EXISTS public.events;
DROP TABLE IF EXISTS public.energy_patterns;
DROP TABLE IF EXISTS public.data_exports;
DROP TABLE IF EXISTS public.daily_metrics;
DROP TABLE IF EXISTS public.communication_methods;
DROP TABLE IF EXISTS public.categories;
DROP TABLE IF EXISTS public.auth_providers;
DROP TABLE IF EXISTS public.audit_logs;
DROP TABLE IF EXISTS public.ai_analysis_jobs;
DROP FUNCTION IF EXISTS public.update_updated_at_column();
DROP FUNCTION IF EXISTS public.update_person_last_interaction();
DROP FUNCTION IF EXISTS public.update_ai_tables_updated_at();
DROP FUNCTION IF EXISTS public.log_audit_event(p_user_id uuid, p_action character varying, p_entity_type character varying, p_entity_id character varying, p_changes jsonb, p_ip_address character varying, p_user_agent text, p_request_id character varying, p_session_id character varying, p_result character varying, p_error_message text);
DROP FUNCTION IF EXISTS public.aggregate_daily_metrics(p_user_id uuid, p_date date);
DROP EXTENSION IF EXISTS "uuid-ossp";
--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: aggregate_daily_metrics(uuid, date); Type: FUNCTION; Schema: public; Owner: vyve
--

CREATE FUNCTION public.aggregate_daily_metrics(p_user_id uuid, p_date date) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    INSERT INTO daily_metrics (
        user_id, date,
        interactions_count,
        unique_persons_count,
        avg_energy_score,
        reflection_completed,
        nudges_generated,
        nudges_acted_on,
        positive_interactions,
        negative_interactions,
        relationships_active
    )
    SELECT 
        p_user_id,
        p_date,
        COUNT(DISTINCT i.id),
        COUNT(DISTINCT i.person_id),
        AVG(CASE 
            WHEN i.energy_impact = 'energizing' THEN 100
            WHEN i.energy_impact = 'neutral' THEN 50
            WHEN i.energy_impact = 'draining' THEN 0
            ELSE 50
        END),
        EXISTS(SELECT 1 FROM reflections WHERE user_id = p_user_id AND DATE(completed_at) = p_date),
        (SELECT COUNT(*) FROM nudges WHERE user_id = p_user_id AND DATE(created_at) = p_date),
        (SELECT COUNT(*) FROM nudges WHERE user_id = p_user_id AND DATE(acted_at) = p_date),
        COUNT(CASE WHEN i.energy_impact = 'energizing' THEN 1 END),
        COUNT(CASE WHEN i.energy_impact = 'draining' THEN 1 END),
        (SELECT COUNT(DISTINCT person_id) FROM interactions 
         WHERE user_id = p_user_id 
         AND interaction_at >= p_date - INTERVAL '30 days'
         AND interaction_at <= p_date)
    FROM interactions i
    WHERE i.user_id = p_user_id
    AND DATE(i.interaction_at) = p_date
    GROUP BY p_user_id, p_date
    ON CONFLICT (user_id, date) DO UPDATE SET
        interactions_count = EXCLUDED.interactions_count,
        unique_persons_count = EXCLUDED.unique_persons_count,
        avg_energy_score = EXCLUDED.avg_energy_score,
        reflection_completed = EXCLUDED.reflection_completed,
        nudges_generated = EXCLUDED.nudges_generated,
        nudges_acted_on = EXCLUDED.nudges_acted_on,
        positive_interactions = EXCLUDED.positive_interactions,
        negative_interactions = EXCLUDED.negative_interactions,
        relationships_active = EXCLUDED.relationships_active,
        updated_at = NOW();
END;
$$;


ALTER FUNCTION public.aggregate_daily_metrics(p_user_id uuid, p_date date) OWNER TO vyve;

--
-- Name: log_audit_event(uuid, character varying, character varying, character varying, jsonb, character varying, text, character varying, character varying, character varying, text); Type: FUNCTION; Schema: public; Owner: vyve
--

CREATE FUNCTION public.log_audit_event(p_user_id uuid, p_action character varying, p_entity_type character varying, p_entity_id character varying, p_changes jsonb, p_ip_address character varying, p_user_agent text, p_request_id character varying, p_session_id character varying, p_result character varying, p_error_message text DEFAULT NULL::text) RETURNS uuid
    LANGUAGE plpgsql
    AS $$
DECLARE
    audit_id UUID;
BEGIN
    INSERT INTO audit_logs (
        user_id, action, entity_type, entity_id, changes,
        ip_address, user_agent, request_id, session_id,
        result, error_message
    ) VALUES (
        p_user_id, p_action, p_entity_type, p_entity_id, p_changes,
        p_ip_address, p_user_agent, p_request_id, p_session_id,
        p_result, p_error_message
    ) RETURNING id INTO audit_id;
    
    RETURN audit_id;
END;
$$;


ALTER FUNCTION public.log_audit_event(p_user_id uuid, p_action character varying, p_entity_type character varying, p_entity_id character varying, p_changes jsonb, p_ip_address character varying, p_user_agent text, p_request_id character varying, p_session_id character varying, p_result character varying, p_error_message text) OWNER TO vyve;

--
-- Name: update_ai_tables_updated_at(); Type: FUNCTION; Schema: public; Owner: vyve
--

CREATE FUNCTION public.update_ai_tables_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_ai_tables_updated_at() OWNER TO vyve;

--
-- Name: update_person_last_interaction(); Type: FUNCTION; Schema: public; Owner: vyve
--

CREATE FUNCTION public.update_person_last_interaction() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE people 
    SET 
        last_interaction_at = NEW.interaction_at,
        interaction_count = interaction_count + 1
    WHERE id = NEW.person_id;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_person_last_interaction() OWNER TO vyve;

--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: vyve
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO vyve;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: ai_analysis_jobs; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.ai_analysis_jobs (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    job_type character varying(50) NOT NULL,
    status character varying(50) DEFAULT 'pending'::character varying NOT NULL,
    priority integer DEFAULT 5,
    person_ids text[],
    total_items integer DEFAULT 0,
    processed_items integer DEFAULT 0,
    failed_items integer DEFAULT 0,
    progress numeric(5,2) DEFAULT 0,
    result_data jsonb,
    error text,
    total_tokens_used integer DEFAULT 0,
    estimated_cost numeric(10,4) DEFAULT 0,
    started_at timestamp without time zone,
    completed_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.ai_analysis_jobs OWNER TO vyve;

--
-- Name: audit_logs; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.audit_logs (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    action character varying(100) NOT NULL,
    entity_type character varying(50),
    entity_id character varying(255),
    changes jsonb,
    ip_address character varying(45),
    user_agent text,
    request_id character varying(255),
    session_id character varying(255),
    result character varying(20),
    error_message text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.audit_logs OWNER TO vyve;

--
-- Name: auth_providers; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.auth_providers (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    provider character varying(50) NOT NULL,
    provider_id character varying(255) NOT NULL,
    access_token text,
    refresh_token text,
    expires_at timestamp without time zone,
    raw_data jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.auth_providers OWNER TO vyve;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.categories (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    name character varying(100) NOT NULL,
    color character varying(20),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.categories OWNER TO vyve;

--
-- Name: communication_methods; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.communication_methods (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(50) NOT NULL,
    icon character varying(50),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.communication_methods OWNER TO vyve;

--
-- Name: daily_metrics; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.daily_metrics (
    user_id uuid NOT NULL,
    date date NOT NULL,
    interactions_count integer DEFAULT 0,
    unique_persons_count integer DEFAULT 0,
    avg_energy_score double precision DEFAULT 0,
    reflection_completed boolean DEFAULT false,
    nudges_generated integer DEFAULT 0,
    nudges_acted_on integer DEFAULT 0,
    positive_interactions integer DEFAULT 0,
    negative_interactions integer DEFAULT 0,
    relationships_active integer DEFAULT 0,
    relationships_improved integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.daily_metrics OWNER TO vyve;

--
-- Name: data_exports; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.data_exports (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    format character varying(10) DEFAULT 'json'::character varying,
    file_url text,
    file_size bigint,
    requested_at timestamp without time zone DEFAULT now() NOT NULL,
    completed_at timestamp without time zone,
    expires_at timestamp without time zone,
    error text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.data_exports OWNER TO vyve;

--
-- Name: energy_patterns; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.energy_patterns (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(50) NOT NULL,
    color character varying(20),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.energy_patterns OWNER TO vyve;

--
-- Name: events; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.events (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    event_type character varying(100) NOT NULL,
    properties jsonb DEFAULT '{}'::jsonb,
    session_id character varying(255),
    ip_address character varying(45),
    user_agent text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.events OWNER TO vyve;

--
-- Name: intentions; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.intentions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(50) NOT NULL,
    color character varying(20),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.intentions OWNER TO vyve;

--
-- Name: interactions; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.interactions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    person_id uuid NOT NULL,
    energy_impact character varying(20) NOT NULL,
    context text[],
    duration integer,
    quality integer,
    notes text,
    location character varying(255),
    special_tags text[],
    interaction_at timestamp without time zone DEFAULT now() NOT NULL,
    metadata jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.interactions OWNER TO vyve;

--
-- Name: nudges; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.nudges (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    person_id uuid,
    type character varying(50) NOT NULL,
    title character varying(255) NOT NULL,
    message text NOT NULL,
    priority character varying(20) DEFAULT 'medium'::character varying,
    action character varying(100),
    action_data jsonb,
    seen boolean DEFAULT false,
    acted_on boolean DEFAULT false,
    seen_at timestamp without time zone,
    acted_at timestamp without time zone,
    expires_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    source character varying(20) DEFAULT 'system'::character varying,
    reasoning text,
    suggested_actions text[],
    conversation_starters text[],
    timing character varying(50),
    estimated_impact character varying(20),
    status character varying(50) DEFAULT 'pending'::character varying,
    accepted_at timestamp without time zone,
    completed_at timestamp without time zone,
    dismissed_at timestamp without time zone,
    provider character varying(50),
    model character varying(100),
    analysis_id uuid
);


ALTER TABLE public.nudges OWNER TO vyve;

--
-- Name: COLUMN nudges.source; Type: COMMENT; Schema: public; Owner: vyve
--

COMMENT ON COLUMN public.nudges.source IS 'Source of the nudge: ai (AI-generated) or system (rule-based)';


--
-- Name: COLUMN nudges.status; Type: COMMENT; Schema: public; Owner: vyve
--

COMMENT ON COLUMN public.nudges.status IS 'Current status: pending, seen, accepted, completed, dismissed';


--
-- Name: COLUMN nudges.analysis_id; Type: COMMENT; Schema: public; Owner: vyve
--

COMMENT ON COLUMN public.nudges.analysis_id IS 'Links to relationship_analyses if this is an AI-generated nudge';


--
-- Name: people; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.people (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    name character varying(255) NOT NULL,
    category_id uuid,
    relationship character varying(100),
    avatar_url text,
    health_score double precision DEFAULT 50.0,
    energy_pattern_id uuid,
    last_interaction_at timestamp without time zone,
    interaction_count integer DEFAULT 0,
    communication_method_id uuid,
    relationship_status_id uuid,
    intention_id uuid,
    context text[],
    notes text,
    custom_fields jsonb DEFAULT '{}'::jsonb,
    reminder_frequency character varying(20),
    next_reminder_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.people OWNER TO vyve;

--
-- Name: push_tokens; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.push_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token text NOT NULL,
    platform character varying(20) NOT NULL,
    device_id character varying(255),
    device_info jsonb,
    active boolean DEFAULT true,
    last_used_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.push_tokens OWNER TO vyve;

--
-- Name: reflections; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.reflections (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    prompt text,
    responses text[],
    mood character varying(50),
    energy_level integer,
    insights text[],
    intentions text[],
    gratitude text[],
    metadata jsonb DEFAULT '{}'::jsonb,
    completed_at timestamp without time zone DEFAULT now() NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    CONSTRAINT reflections_energy_level_check CHECK (((energy_level >= 1) AND (energy_level <= 10)))
);


ALTER TABLE public.reflections OWNER TO vyve;

--
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.refresh_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token text NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    revoked boolean DEFAULT false,
    revoked_at timestamp without time zone,
    ip_address character varying(45),
    user_agent text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.refresh_tokens OWNER TO vyve;

--
-- Name: relationship_analyses; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.relationship_analyses (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    person_id uuid NOT NULL,
    connection_strength numeric(5,2) DEFAULT 0,
    engagement_quality numeric(5,2) DEFAULT 0,
    communication_balance numeric(5,2) DEFAULT 0,
    energy_alignment numeric(5,2) DEFAULT 0,
    relationship_health numeric(5,2) DEFAULT 0,
    overall_score numeric(5,2) DEFAULT 0,
    summary text,
    key_insights text[],
    patterns text[],
    strengths text[],
    concerns text[],
    trend_direction character varying(50),
    provider character varying(50) NOT NULL,
    model character varying(100),
    tokens_used integer DEFAULT 0,
    processing_time_ms integer DEFAULT 0,
    version integer DEFAULT 1,
    analyzed_at timestamp without time zone DEFAULT now() NOT NULL,
    interactions_count integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.relationship_analyses OWNER TO vyve;

--
-- Name: relationship_statuses; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.relationship_statuses (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(50) NOT NULL,
    color character varying(20),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.relationship_statuses OWNER TO vyve;

--
-- Name: user_consents; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.user_consents (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    consent_type character varying(50) NOT NULL,
    granted boolean NOT NULL,
    version character varying(20),
    ip_address character varying(45),
    user_agent text,
    granted_at timestamp without time zone NOT NULL,
    revoked_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.user_consents OWNER TO vyve;

--
-- Name: users; Type: TABLE; Schema: public; Owner: vyve
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username text NOT NULL,
    email character varying(255) NOT NULL,
    email_verified boolean DEFAULT false,
    password_hash text,
    avatar_url text,
    display_name character varying(100),
    bio text,
    timezone character varying(50) DEFAULT 'UTC'::character varying,
    locale character varying(10) DEFAULT 'en'::character varying,
    last_login_at timestamp without time zone,
    last_activity_at timestamp without time zone,
    streak_count integer DEFAULT 0,
    last_reflection_at timestamp without time zone,
    settings jsonb DEFAULT '{}'::jsonb,
    metadata jsonb DEFAULT '{}'::jsonb,
    data_residency character varying(10) DEFAULT 'us'::character varying,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    onboarding_completed boolean DEFAULT false,
    onboarding_steps jsonb DEFAULT '[]'::jsonb
);


ALTER TABLE public.users OWNER TO vyve;

--
-- Data for Name: ai_analysis_jobs; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.ai_analysis_jobs (id, user_id, job_type, status, priority, person_ids, total_items, processed_items, failed_items, progress, result_data, error, total_tokens_used, estimated_cost, started_at, completed_at, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: audit_logs; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.audit_logs (id, user_id, action, entity_type, entity_id, changes, ip_address, user_agent, request_id, session_id, result, error_message, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: auth_providers; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.auth_providers (id, user_id, provider, provider_id, access_token, refresh_token, expires_at, raw_data, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.categories (id, user_id, name, color, created_at, updated_at, deleted_at) FROM stdin;
d23e6f34-3ed3-4453-94af-e4bc10b91b1f	0e1a954b-851b-493a-8a64-753b15e20a70	Friend	\N	2025-09-10 12:28:52.034306	2025-09-10 12:28:52.034306	\N
f73bf100-378d-4991-8e27-ec1facc90515	0e1a954b-851b-493a-8a64-753b15e20a70	Family	\N	2025-09-10 12:29:18.584611	2025-09-10 12:29:18.584611	\N
23eb167c-59a5-42c9-9038-61b30fc3d48a	0e1a954b-851b-493a-8a64-753b15e20a70	Colleague	\N	2025-09-10 12:29:55.371462	2025-09-10 12:29:55.371462	\N
81a0589c-5a84-4ba9-a832-b55db5ee0527	0e1a954b-851b-493a-8a64-753b15e20a70	Collaborator	\N	2025-09-10 12:30:48.8637	2025-09-10 12:30:48.8637	\N
0fc4a313-baf7-4f93-a91b-17e178aa9ece	0e1a954b-851b-493a-8a64-753b15e20a70	Client	\N	2025-09-10 12:31:21.214333	2025-09-10 12:31:21.214333	\N
ec22df8d-f4a8-4ae3-8c4c-dc4aa9246a24	0e1a954b-851b-493a-8a64-753b15e20a70	Boss	\N	2025-09-10 12:31:32.245412	2025-09-10 12:31:32.245412	\N
77a4dac9-e3c6-41d6-9e35-56b5604010d6	0e1a954b-851b-493a-8a64-753b15e20a70	Acquaintance	\N	2025-09-10 12:32:12.671727	2025-09-10 12:32:12.671727	\N
\.


--
-- Data for Name: communication_methods; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.communication_methods (id, name, icon, created_at, updated_at, deleted_at) FROM stdin;
90635e6d-5539-487a-8fb0-2b9374d255cd	whatsapp	whatsapp	2025-08-14 10:43:58.885299	2025-08-14 10:43:58.885299	\N
5d5c54ba-cdb0-445a-80a9-61d3e7884fcf	call	phone	2025-08-14 10:43:58.885299	2025-08-14 10:43:58.885299	\N
7a534d39-e4cf-4b5b-99c1-7e77b1e92db7	text	sms	2025-08-14 10:43:58.885299	2025-08-14 10:43:58.885299	\N
5394a706-1d0d-464b-a9ff-531281661585	email	email	2025-08-14 10:43:58.885299	2025-08-14 10:43:58.885299	\N
8ba6023c-f38b-47d6-a40f-35e9fedc0262	slack	slack	2025-09-12 10:35:53.660955	2025-09-12 10:35:53.660955	\N
27ad1704-41a3-4b89-b682-bac4b8d05839	in person	in_person	2025-09-12 10:36:16.247324	2025-09-12 10:36:16.247324	\N
\.


--
-- Data for Name: daily_metrics; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.daily_metrics (user_id, date, interactions_count, unique_persons_count, avg_energy_score, reflection_completed, nudges_generated, nudges_acted_on, positive_interactions, negative_interactions, relationships_active, relationships_improved, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: data_exports; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.data_exports (id, user_id, status, format, file_url, file_size, requested_at, completed_at, expires_at, error, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: energy_patterns; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.energy_patterns (id, name, color, created_at, updated_at, deleted_at) FROM stdin;
c596933f-845c-4d0d-94ab-e262deded367	energizing	green	2025-08-14 10:43:58.854461	2025-08-14 10:43:58.854461	\N
b34480fa-eea5-46b0-9a94-d2bb7a01eaae	neutral	gray	2025-08-14 10:43:58.854461	2025-08-14 10:43:58.854461	\N
623e8c45-064d-46fc-af11-281c979616cf	draining	red	2025-08-14 10:43:58.854461	2025-08-14 10:43:58.854461	\N
\.


--
-- Data for Name: events; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.events (id, user_id, event_type, properties, session_id, ip_address, user_agent, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: intentions; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.intentions (id, name, color, created_at, updated_at, deleted_at) FROM stdin;
18848a05-94d3-4c04-8126-ad48ae1a3196	audit	gray	2025-08-14 10:43:58.92776	2025-08-14 10:43:58.92776	\N
3381e24a-38a1-4f7d-86f7-e21df8bd9414	improve	green	2025-08-14 10:43:58.92776	2025-08-14 10:43:58.92776	\N
6dfc9d23-c815-4567-bdb9-2fb13045dc47	maintain	blue	2025-08-14 10:43:58.92776	2025-08-14 10:43:58.92776	\N
a8102418-0770-43cd-989d-e27049cf6b6c	boundaries	red	2025-08-14 10:43:58.92776	2025-08-14 10:43:58.92776	\N
\.


--
-- Data for Name: interactions; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.interactions (id, user_id, person_id, energy_impact, context, duration, quality, notes, location, special_tags, interaction_at, metadata, created_at, updated_at, deleted_at) FROM stdin;
fe9a88b7-76e5-4820-9cc0-774f4970d9b7	0e1a954b-851b-493a-8a64-753b15e20a70	d68e2355-581f-4476-b1e3-4a85e03a182e	neutral	{}	0	0			{}	2025-09-10 23:35:59.349583	\N	2025-09-10 23:35:59.358217	2025-09-10 23:35:59.358217	\N
3bc971ba-236c-433e-9f20-65b1954a372c	0e1a954b-851b-493a-8a64-753b15e20a70	b82af2ee-5db4-4103-8bf6-fbb8cca8bd07	neutral	{}	0	0			{}	2025-09-12 10:12:34.373198	\N	2025-09-12 10:12:34.37503	2025-09-12 10:12:34.37503	\N
01ee8272-9224-4ad7-bb71-42005f52e9f2	0e1a954b-851b-493a-8a64-753b15e20a70	2e76ea47-c7f9-40f5-a369-0e894770808b	neutral	{}	0	0			{}	2025-09-12 10:40:09.594335	\N	2025-09-12 10:40:09.595836	2025-09-12 10:40:09.595836	\N
5a4e6fdb-1299-4075-b700-d6790a70485a	0e1a954b-851b-493a-8a64-753b15e20a70	2e76ea47-c7f9-40f5-a369-0e894770808b	neutral	{}	0	0			{}	2025-09-12 10:50:40.807451	\N	2025-09-12 10:50:40.809726	2025-09-12 10:50:40.809726	\N
094f8d26-bdaa-4566-8828-a4f3863dc6c9	0e1a954b-851b-493a-8a64-753b15e20a70	9b476068-66aa-44ca-8325-fe80910143ca	neutral	{}	0	0			{}	2025-09-12 23:23:24.982874	\N	2025-09-12 23:23:24.983925	2025-09-12 23:23:24.983925	\N
1e7abbd9-894f-4188-b4a2-29c349c05590	0e1a954b-851b-493a-8a64-753b15e20a70	b82af2ee-5db4-4103-8bf6-fbb8cca8bd07	neutral	{}	0	0	Call interaction via map		{}	2025-09-29 19:46:45.785155	\N	2025-09-29 19:46:45.787977	2025-09-29 19:46:45.787977	\N
187b560f-3918-4f17-b424-7dfb61f4d338	0e1a954b-851b-493a-8a64-753b15e20a70	9b476068-66aa-44ca-8325-fe80910143ca	neutral	{}	0	0	Email interaction via map		{}	2025-09-29 20:02:10.445979	\N	2025-09-29 20:02:10.447596	2025-09-29 20:02:10.447596	\N
1e1eac0f-a716-4e24-9a42-eea206c67d6c	0e1a954b-851b-493a-8a64-753b15e20a70	9b476068-66aa-44ca-8325-fe80910143ca	energizing	{}	0	0	energizing interaction via map		{}	2025-09-29 20:22:55.015148	\N	2025-09-29 20:22:55.016993	2025-09-29 20:22:55.016993	\N
f2f313ee-2c22-4230-936f-0ccfff27165f	0e1a954b-851b-493a-8a64-753b15e20a70	d68e2355-581f-4476-b1e3-4a85e03a182e	draining	{}	0	0	draining interaction via map		{}	2025-09-29 20:23:08.302483	\N	2025-09-29 20:23:08.308929	2025-09-29 20:23:08.308929	\N
bd877111-ea28-4fa7-a852-c3a94fa1a197	0e1a954b-851b-493a-8a64-753b15e20a70	b82af2ee-5db4-4103-8bf6-fbb8cca8bd07	energizing	{}	0	0	energizing interaction via map		{}	2025-09-30 12:58:30.984236	\N	2025-09-30 12:58:30.998057	2025-09-30 12:58:30.998057	\N
bc8937aa-c651-4600-8289-963423fe856d	0e1a954b-851b-493a-8a64-753b15e20a70	2e76ea47-c7f9-40f5-a369-0e894770808b	energizing	{}	0	0	energizing interaction via map		{}	2025-09-30 12:58:56.803042	\N	2025-09-30 12:58:56.804234	2025-09-30 12:58:56.804234	\N
7630072c-399a-4ae3-b200-6d1f3e2f36c4	0e1a954b-851b-493a-8a64-753b15e20a70	b82af2ee-5db4-4103-8bf6-fbb8cca8bd07	neutral	{}	0	0			{}	2025-09-30 13:15:22.7371	\N	2025-09-30 13:15:22.738015	2025-09-30 13:15:22.738015	\N
95a109f6-ac78-4af8-8a71-9dc06beaa108	0e1a954b-851b-493a-8a64-753b15e20a70	2e76ea47-c7f9-40f5-a369-0e894770808b	energizing	{}	0	0	energizing interaction via map		{}	2025-09-30 16:48:35.298369	\N	2025-09-30 16:48:35.307833	2025-09-30 16:48:35.307833	\N
b08a9dde-5897-4775-b243-e459f12aa99b	0e1a954b-851b-493a-8a64-753b15e20a70	9b476068-66aa-44ca-8325-fe80910143ca	neutral	{}	0	0	neutral interaction via map		{}	2025-09-30 16:48:54.385065	\N	2025-09-30 16:48:54.385781	2025-09-30 16:48:54.385781	\N
2d0b26d0-65b0-48d7-a8c8-42e8b2a7af46	0e1a954b-851b-493a-8a64-753b15e20a70	3a994d2f-aaba-41e9-999a-0fcb81ff8ff4	energizing	{}	0	0	energizing interaction via map		{}	2025-09-30 16:49:02.046377	\N	2025-09-30 16:49:02.047797	2025-09-30 16:49:02.047797	\N
2b858022-2230-4cb4-93f3-3a61b486bea3	0e1a954b-851b-493a-8a64-753b15e20a70	da7356bc-8bcd-4084-9a45-1c84d3cb2d84	energizing	{}	0	0	energizing interaction via map		{}	2025-09-30 16:50:39.759492	\N	2025-09-30 16:50:39.760761	2025-09-30 16:50:39.760761	\N
f8e81d46-1cd1-484c-82d6-b0b689360b35	0e1a954b-851b-493a-8a64-753b15e20a70	b82af2ee-5db4-4103-8bf6-fbb8cca8bd07	neutral	{}	0	0	neutral interaction via map		{}	2025-09-30 16:56:54.232901	\N	2025-09-30 16:56:54.24055	2025-09-30 16:56:54.24055	\N
c6662176-a717-4fa0-8dec-6194fe33e9c2	0e1a954b-851b-493a-8a64-753b15e20a70	da7356bc-8bcd-4084-9a45-1c84d3cb2d84	energizing	{}	0	0	energizing interaction via map		{}	2025-09-30 17:31:11.283867	\N	2025-09-30 17:31:11.286015	2025-09-30 17:31:11.286015	\N
4502c4c0-8026-4ca3-8f18-3cd524a07c3f	0e1a954b-851b-493a-8a64-753b15e20a70	da7356bc-8bcd-4084-9a45-1c84d3cb2d84	draining	{}	0	0	draining interaction via map		{}	2025-09-30 17:31:14.318903	\N	2025-09-30 17:31:14.319413	2025-09-30 17:31:14.319413	\N
ae8e6731-2a26-4fc9-a3bc-60efb6109d39	0e1a954b-851b-493a-8a64-753b15e20a70	d68e2355-581f-4476-b1e3-4a85e03a182e	draining	{}	0	0	draining interaction via map		{}	2025-09-30 17:34:40.990757	\N	2025-09-30 17:34:40.992618	2025-09-30 17:34:40.992618	\N
01f99acc-7377-4f76-a050-46ea5c030a6f	0e1a954b-851b-493a-8a64-753b15e20a70	b82af2ee-5db4-4103-8bf6-fbb8cca8bd07	draining	{}	0	0	draining interaction via map		{}	2025-09-30 18:27:35.217569	\N	2025-09-30 18:27:35.225038	2025-09-30 18:27:35.225038	\N
5cea4218-f407-4017-ae07-6c248fd65d5e	0e1a954b-851b-493a-8a64-753b15e20a70	d68e2355-581f-4476-b1e3-4a85e03a182e	energizing	{}	0	0	energizing interaction via map		{}	2025-09-30 18:32:13.066212	\N	2025-09-30 18:32:13.067496	2025-09-30 18:32:13.067496	\N
0ab37a31-c79f-465c-968e-b908c3ffd792	0e1a954b-851b-493a-8a64-753b15e20a70	d68e2355-581f-4476-b1e3-4a85e03a182e	energizing	{}	0	0	energizing interaction via map		{}	2025-09-30 18:32:17.890581	\N	2025-09-30 18:32:17.891602	2025-09-30 18:32:17.891602	\N
5b0dd20c-1105-4da4-8d35-7a01baa9b70f	0e1a954b-851b-493a-8a64-753b15e20a70	da7356bc-8bcd-4084-9a45-1c84d3cb2d84	energizing	{}	0	0	energizing interaction via map		{}	2025-09-30 18:48:32.713988	\N	2025-09-30 18:48:32.715886	2025-09-30 18:48:32.715886	\N
9956739c-4e5f-475c-9e3b-6069bf35709b	0e1a954b-851b-493a-8a64-753b15e20a70	3a994d2f-aaba-41e9-999a-0fcb81ff8ff4	energizing	{}	0	0	energizing interaction via map		{}	2025-09-30 18:48:37.643632	\N	2025-09-30 18:48:37.644024	2025-09-30 18:48:37.644024	\N
1096a3ee-4f2f-4083-9dee-913a79317e71	0e1a954b-851b-493a-8a64-753b15e20a70	2e76ea47-c7f9-40f5-a369-0e894770808b	draining	{}	0	0	draining interaction via map		{}	2025-09-30 18:48:40.994848	\N	2025-09-30 18:48:40.998499	2025-09-30 18:48:40.998499	\N
ef8a6bb5-9b15-47cc-b540-acdc3e1138d2	0e1a954b-851b-493a-8a64-753b15e20a70	b82af2ee-5db4-4103-8bf6-fbb8cca8bd07	neutral	{}	0	0	neutral interaction via map		{}	2025-09-30 18:48:44.027004	\N	2025-09-30 18:48:44.027963	2025-09-30 18:48:44.027963	\N
d3d492ca-d350-4043-a5e5-266e6ca5c90f	0e1a954b-851b-493a-8a64-753b15e20a70	b82af2ee-5db4-4103-8bf6-fbb8cca8bd07	neutral	{}	0	0	neutral interaction via map		{}	2025-09-30 18:58:20.765673	\N	2025-09-30 18:58:20.76623	2025-09-30 18:58:20.76623	\N
b732d461-71f8-4690-82ab-1055d4c504db	0e1a954b-851b-493a-8a64-753b15e20a70	9b476068-66aa-44ca-8325-fe80910143ca	energizing	{}	0	0	energizing interaction via map		{}	2025-09-30 18:58:23.930712	\N	2025-09-30 18:58:23.932117	2025-09-30 18:58:23.932117	\N
ba2a6eef-514a-46fe-82d6-bbbe18e64466	0e1a954b-851b-493a-8a64-753b15e20a70	3a994d2f-aaba-41e9-999a-0fcb81ff8ff4	draining	{}	0	0	draining interaction via map		{}	2025-09-30 18:58:27.706825	\N	2025-09-30 18:58:27.708106	2025-09-30 18:58:27.708106	\N
4601e47e-cd5d-4ee7-a8b8-1140abcda3e9	0e1a954b-851b-493a-8a64-753b15e20a70	3a994d2f-aaba-41e9-999a-0fcb81ff8ff4	energizing	{}	0	0	energizing interaction via map		{}	2025-09-30 18:59:15.187456	\N	2025-09-30 18:59:15.189006	2025-09-30 18:59:15.189006	\N
c43360ed-69d7-4027-a48c-d8dddebb6043	0e1a954b-851b-493a-8a64-753b15e20a70	d68e2355-581f-4476-b1e3-4a85e03a182e	energizing	{}	0	0	energizing interaction via map		{}	2025-09-30 18:59:26.737325	\N	2025-09-30 18:59:26.740368	2025-09-30 18:59:26.740368	\N
d540b31e-22f0-4a15-a956-0d99fd52b7c0	0e1a954b-851b-493a-8a64-753b15e20a70	3a994d2f-aaba-41e9-999a-0fcb81ff8ff4	draining	{}	0	0	draining interaction via map		{}	2025-09-30 18:59:49.833924	\N	2025-09-30 18:59:49.836856	2025-09-30 18:59:49.836856	\N
d6c440ee-9190-4ae0-aa13-b0edbe5d9f84	0e1a954b-851b-493a-8a64-753b15e20a70	3a994d2f-aaba-41e9-999a-0fcb81ff8ff4	draining	{}	0	0	draining interaction via map		{}	2025-09-30 18:59:52.458048	\N	2025-09-30 18:59:52.462069	2025-09-30 18:59:52.462069	\N
80fe57b6-e612-470f-8c00-d7732b3a4e06	0e1a954b-851b-493a-8a64-753b15e20a70	9b476068-66aa-44ca-8325-fe80910143ca	draining	{}	0	0	draining interaction via map		{}	2025-10-21 09:26:16.900394	\N	2025-10-21 09:26:16.909247	2025-10-21 09:26:16.909247	\N
8ebbec26-96c2-408f-88bd-31f24a6789ea	0e1a954b-851b-493a-8a64-753b15e20a70	da7356bc-8bcd-4084-9a45-1c84d3cb2d84	energizing	{}	0	0	energizing interaction via map		{}	2025-10-21 09:26:35.023915	\N	2025-10-21 09:26:35.025033	2025-10-21 09:26:35.025033	\N
\.


--
-- Data for Name: nudges; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.nudges (id, user_id, person_id, type, title, message, priority, action, action_data, seen, acted_on, seen_at, acted_at, expires_at, created_at, updated_at, deleted_at, source, reasoning, suggested_actions, conversation_starters, timing, estimated_impact, status, accepted_at, completed_at, dismissed_at, provider, model, analysis_id) FROM stdin;
e74b06b9-1694-4d2d-b6a2-6bcdb39e3880	0e1a954b-851b-493a-8a64-753b15e20a70	9b476068-66aa-44ca-8325-fe80910143ca	reach_out	Initiate Open Dialogue	Start a conversation with family members to express your feelings and challenges. Aim to understand their perspectives and share your own with honesty and empathy.	high		\N	f	f	\N	\N	2025-10-30 23:14:12.090394	2025-10-23 23:14:12.127556	2025-10-23 23:14:12.127556	\N	ai	Open communication can resolve misunderstandings and build a foundation for a healthier relationship.	{}	{}	this_week		pending	\N	\N	\N	openai	gpt-4o	59eecb0b-5ab2-4c27-bf9f-b3b7c29f1e48
afa9b6ae-6257-4440-913a-c4341de8242f	0e1a954b-851b-493a-8a64-753b15e20a70	9b476068-66aa-44ca-8325-fe80910143ca	schedule_call	Regular Check-ins	Establish a routine of regular check-ins with family members to maintain connection and monitor relationship dynamics.	medium		\N	f	f	\N	\N	2025-11-22 23:14:12.291783	2025-10-23 23:14:12.297607	2025-10-23 23:14:12.297607	\N	ai	Regular communication can prevent issues from escalating and promote ongoing understanding.	{}	{}	this_month		pending	\N	\N	\N	openai	gpt-4o	59eecb0b-5ab2-4c27-bf9f-b3b7c29f1e48
63316bdd-1e15-45b2-a662-e8879f883d04	0e1a954b-851b-493a-8a64-753b15e20a70	9b476068-66aa-44ca-8325-fe80910143ca	celebrate	Acknowledge Small Wins	Celebrate small positive interactions or improvements in the relationship to build momentum and positivity.	low		\N	f	f	\N	\N	2025-11-22 23:14:12.304562	2025-10-23 23:14:12.30475	2025-10-23 23:14:12.30475	\N	ai	Recognizing and celebrating small wins can reinforce positive behavior and improve overall relationship satisfaction.	{}	{}	this_month		pending	\N	\N	\N	openai	gpt-4o	59eecb0b-5ab2-4c27-bf9f-b3b7c29f1e48
\.


--
-- Data for Name: people; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.people (id, user_id, name, category_id, relationship, avatar_url, health_score, energy_pattern_id, last_interaction_at, interaction_count, communication_method_id, relationship_status_id, intention_id, context, notes, custom_fields, reminder_frequency, next_reminder_at, created_at, updated_at, deleted_at) FROM stdin;
2e76ea47-c7f9-40f5-a369-0e894770808b	0e1a954b-851b-493a-8a64-753b15e20a70	Luan	ec22df8d-f4a8-4ae3-8c4c-dc4aa9246a24			46.35036496350365	\N	2025-09-30 18:48:41.001797	10	8ba6023c-f38b-47d6-a40f-35e9fedc0262	f29de9cd-18ae-42f8-8849-c33ba654bfd5	a8102418-0770-43cd-989d-e27049cf6b6c	{Boss}	You're taking a boundaries approach to this stable boss relationship. This awareness is the first step toward intentional connection.	\N		\N	2025-09-12 10:39:20.058143	2025-09-30 18:48:41.005942	\N
b82af2ee-5db4-4103-8bf6-fbb8cca8bd07	0e1a954b-851b-493a-8a64-753b15e20a70	Elona	\N			46.93385895751205	\N	2025-09-30 18:58:20.774235	16	\N	\N	\N	{Friend}	You're taking a improve / reconnect approach to this stable / solid friend relationship. This awareness is the first step toward intentional connection.	\N		\N	2025-09-08 20:39:24.39873	2025-09-30 18:58:20.780854	\N
d68e2355-581f-4476-b1e3-4a85e03a182e	0e1a954b-851b-493a-8a64-753b15e20a70	Person	\N	friend		78.23129251700682	\N	2025-09-30 18:59:26.742509	12	\N	\N	\N	{}	Created via API test	\N		\N	2025-09-06 22:20:35.361963	2025-09-30 18:59:26.746104	\N
da7356bc-8bcd-4084-9a45-1c84d3cb2d84	0e1a954b-851b-493a-8a64-753b15e20a70	Mark	\N			85.4014598540146	\N	2025-10-21 09:26:35.028566	10	\N	\N	\N	{Friend}	You're taking a improve / reconnect approach to this stable / solid friend relationship. This awareness is the first step toward intentional connection.	\N		\N	2025-09-08 20:39:57.895655	2025-10-21 09:26:35.035632	\N
33828423-9846-4ed9-932b-779dfbf923d5	0e1a954b-851b-493a-8a64-753b15e20a70	Jacob	ec22df8d-f4a8-4ae3-8c4c-dc4aa9246a24			50	\N	\N	0	5d5c54ba-cdb0-445a-80a9-61d3e7884fcf	bc29f97e-2542-437c-a01a-adb8a09d7c39	a8102418-0770-43cd-989d-e27049cf6b6c	{Boss}	You're taking a boundaries approach to this new boss relationship. This awareness is the first step toward intentional connection.	\N		\N	2025-10-21 09:31:45.628916	2025-10-21 09:31:45.628916	\N
3a994d2f-aaba-41e9-999a-0fcb81ff8ff4	0e1a954b-851b-493a-8a64-753b15e20a70	Luan Jubica	77a4dac9-e3c6-41d6-9e35-56b5604010d6			28.571428571428573	\N	2025-09-30 18:59:52.46452	12	90635e6d-5539-487a-8fb0-2b9374d255cd	299c5233-e0ff-4c12-83bd-737add383854	18848a05-94d3-4c04-8126-ad48ae1a3196	{Acquaintance}	You're taking a audit approach to this fading acquaintance relationship. This awareness is the first step toward intentional connection.	\N		\N	2025-09-10 12:33:01.363228	2025-09-30 18:59:52.467126	\N
9b476068-66aa-44ca-8325-fe80910143ca	0e1a954b-851b-493a-8a64-753b15e20a70	Blendina	\N			44.89795918367347	\N	2025-10-21 09:26:16.931143	12	\N	\N	\N	{Family}	You're taking a just audit it approach to this new / early stage family relationship. This awareness is the first step toward intentional connection.	\N		\N	2025-09-08 20:38:00.550258	2025-10-21 09:26:16.936784	\N
\.


--
-- Data for Name: push_tokens; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.push_tokens (id, user_id, token, platform, device_id, device_info, active, last_used_at, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: reflections; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.reflections (id, user_id, prompt, responses, mood, energy_level, insights, intentions, gratitude, metadata, completed_at, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.refresh_tokens (id, user_id, token, expires_at, revoked, revoked_at, ip_address, user_agent, created_at, updated_at, deleted_at) FROM stdin;
4b09e7fb-bbac-4eb9-a663-2c8accf054ad	0e1a954b-851b-493a-8a64-753b15e20a70	BWXuZ4SQELcncfnnyONuPEzp-O0oFShug6LImHF7sqn_1hLDzqyn2z6vTU76pzCJ	2025-08-23 17:39:00.529389	f	\N			2025-08-16 17:39:00.529783	2025-08-16 17:39:00.529783	\N
3db16d2b-201e-4803-b65c-7271c6ecf3b8	0e1a954b-851b-493a-8a64-753b15e20a70	0BgAqvO4kQsEyGJfwNPWqqE9XjUf3er7DLMeOrcXLXtcPWnwldu4XFLLK5T1QWnW	2025-08-23 17:39:34.864469	f	\N			2025-08-16 17:39:34.869012	2025-08-16 17:39:34.869012	\N
02c28b43-00e0-4d0b-83af-7ee2355ac0d7	0e1a954b-851b-493a-8a64-753b15e20a70	Og2et1ounWD7dLupG8o8V5rKD4rJr4sfL_jhvte0xtWyW6w_2lDDuv_XCGeC8CDm	2025-08-23 17:43:40.351613	f	\N			2025-08-16 17:43:40.352142	2025-08-16 17:43:40.352142	\N
fcb6280e-4d1d-47c6-8066-9896da9f9786	0e1a954b-851b-493a-8a64-753b15e20a70	DgbVL9Af6O5q3YHXkXJ6aE4vVcScg-k66gEzweAycYeiui1Drevmbr9FHneP4OIX	2025-08-23 17:45:55.347575	f	\N			2025-08-16 17:45:55.34811	2025-08-16 17:45:55.34811	\N
f49096f4-9d7e-4c99-b156-10b72da50524	0e1a954b-851b-493a-8a64-753b15e20a70	0cvQBas-Wnz3Y-M_8EDlOKwi28llXunuM39c-DMNNVJnHoGPKPBu_nCcnMAzXz6x	2025-08-23 22:05:55.600325	f	\N			2025-08-16 22:05:55.600842	2025-08-16 22:05:55.600842	\N
d4f4019d-a55a-4b8e-84ea-ac17aaad4945	0e1a954b-851b-493a-8a64-753b15e20a70	cBE1cTNu4XPtdKi45d62736nrO7ot92q-RhBEp_NescnAi80fLPo48tHbdBpU5Ma	2025-08-26 11:04:19.735104	f	\N			2025-08-19 11:04:19.736136	2025-08-19 11:04:19.736136	\N
bcc9f399-7e84-4690-ba17-48e5b14d9943	0e1a954b-851b-493a-8a64-753b15e20a70	EQ10_YXs7TiVQsWDw-NY0tzl08-pEcjYBUwuc7LDTPgMzA0ezQabFjiw6ikq6Bzo	2025-08-26 11:41:18.591192	f	\N			2025-08-19 11:41:18.591501	2025-08-19 11:41:18.591501	\N
f0266d1a-0793-4dc5-98d3-9bff6aea29e8	0e1a954b-851b-493a-8a64-753b15e20a70	h3gV5Y2OKPeV4Vj8B6aLrzzegaR-gcAzZ8x57xtb9SCNBSXkFnPCh-aHU-DqmoEu	2025-08-27 21:10:20.178277	f	\N			2025-08-20 21:10:20.17878	2025-08-20 21:10:20.17878	\N
fdf8f79b-974a-4345-b1c7-afbb8ffa841a	0e1a954b-851b-493a-8a64-753b15e20a70	1Vs15sc0AbLBrm0Zj8zOfkMdDpxjdBml1GmN7QVGHmbRUyEFIa4tsXVqVnyoUjuT	2025-08-27 21:11:33.024753	f	\N			2025-08-20 21:11:33.025159	2025-08-20 21:11:33.025159	\N
ef2a5027-6a49-48ce-b780-711096c9a4df	0e1a954b-851b-493a-8a64-753b15e20a70	aa46sAnZMFajkChsHFuP8wI87WbUJvWILDrgmUWgykNYVHHQEiRvj_7Z4zWwlvZT	2025-09-13 17:27:47.146965	t	2025-09-06 17:35:16.707753			2025-09-06 17:27:47.147486	2025-09-06 17:35:16.707799	\N
6114d570-b223-4f26-a358-1a55e88af6b0	0e1a954b-851b-493a-8a64-753b15e20a70	FTfO5ghgWygjw2pXyXhxp81PvH1Kkn2IGxaNieIM_Vs26BcRcQnSoBjj-jIKm6mO	2025-09-13 17:35:16.71154	t	2025-09-06 17:38:19.456526			2025-09-06 17:35:16.711766	2025-09-06 17:38:19.456799	\N
46dc6bbd-4bc8-48c8-9c6e-fa797d086031	0e1a954b-851b-493a-8a64-753b15e20a70	oyGVNFoZ4YtZaPdOoVojPF7xFTHdXRHVqPJIa950b8mMdGf04XbxjLGDF4Bq0mUB	2025-09-13 17:39:10.385316	f	\N			2025-09-06 17:39:10.385449	2025-09-06 17:39:10.385449	\N
42dbf2d9-7eff-4df3-8f88-07c13552f23e	0e1a954b-851b-493a-8a64-753b15e20a70	pDIOtTIWfZJz7TQ1AySE9BAudSp0qG0XCc9MfSIsv_ct-VP7e4qTRI2jPhUDHpdZ	2025-09-13 19:40:14.053863	f	\N			2025-09-06 19:40:14.054375	2025-09-06 19:40:14.054375	\N
4c39044f-cf3e-43c1-be25-b09bb977fd93	0e1a954b-851b-493a-8a64-753b15e20a70	iqWGkRSMH39t6_8hG9AZNC9Q44eFMhHbiu23cXke7IH8ZQaJ2YvPBExlP3H6cR8n	2025-09-13 20:33:59.854737	f	\N			2025-09-06 20:33:59.857134	2025-09-06 20:33:59.857134	\N
4ee092bf-9a51-44e1-9445-307b8d153825	0e1a954b-851b-493a-8a64-753b15e20a70	edZTidWzmKEDfN0IlxHXGryt4E0hmvHjk2HdscYSdsm6OpgkFpjNIIWGGONdkiaW	2025-09-13 21:36:05.708755	f	\N			2025-09-06 21:36:05.709671	2025-09-06 21:36:05.709671	\N
04aef08c-6ba6-4d2f-8a2b-1d5d8437715a	0e1a954b-851b-493a-8a64-753b15e20a70	W0ma9DW38x-7IW8u3QtXCKY7vOzsRjdDt02WrTHNRVEIMOva_fiXzWDN6pDo19ee	2025-09-13 22:20:35.198708	f	\N			2025-09-06 22:20:35.199269	2025-09-06 22:20:35.199269	\N
1f9bc1e8-d765-4abc-a3a3-cfb74f1ccfba	0e1a954b-851b-493a-8a64-753b15e20a70	yxnZ91X_touL4IkQco38VD9adwpeYNhLoDIpSdxNpsbgjVTfpo0ZBwRnyHPoKe-6	2025-09-13 22:23:41.356448	f	\N			2025-09-06 22:23:41.356918	2025-09-06 22:23:41.356918	\N
2ce14224-db0f-48e4-b7b1-d2285a404777	0e1a954b-851b-493a-8a64-753b15e20a70	4PgyUUiWW_U4YXYUcFUIeJ6pieJElHwqSlHdLuD39eHmd51gwq48WkxLXDW6eTay	2025-09-13 22:26:05.385444	f	\N			2025-09-06 22:26:05.385806	2025-09-06 22:26:05.385806	\N
18428257-c995-4e89-a1d2-5214a92542ba	0e1a954b-851b-493a-8a64-753b15e20a70	oV7gOoV43qa7PumaHD_fmf3azYmZevA22Yy61qQfmzxwEHtUV-B32EEd-y_2BwCf	2025-09-13 22:34:12.404774	f	\N			2025-09-06 22:34:12.405068	2025-09-06 22:34:12.405068	\N
61b6aeac-c14d-40ac-9a8b-854bb2a063a5	0e1a954b-851b-493a-8a64-753b15e20a70	LuyfvANKdVVEr8NedHmqnIXSMV13yk2ju0aboVulmrxuIjdDPG0mCyOopp8oxTqs	2025-09-13 22:47:56.664668	f	\N			2025-09-06 22:47:56.665065	2025-09-06 22:47:56.665065	\N
fb114cd4-c5da-49e4-883c-726a527220b6	0e1a954b-851b-493a-8a64-753b15e20a70	4cj23ZRBxAdUaZJdsIBg-LYns9EGUgCcFyqSdNBVA4zr3n_BvG4b_XGaqfv94rTv	2025-09-13 22:48:55.312591	f	\N			2025-09-06 22:48:55.312743	2025-09-06 22:48:55.312743	\N
134042c4-07eb-45df-bdbe-accdf420a3b6	0e1a954b-851b-493a-8a64-753b15e20a70	LFuicVISUsj2p5UwVUHnVdX95qyGr9H5TbmsB9p0nO2dRmkJDGf4-e3srxJs3V0p	2025-09-13 23:06:13.786211	f	\N			2025-09-06 23:06:13.787091	2025-09-06 23:06:13.787091	\N
0e0961b1-96f0-4e7f-a98a-37bdadfcfc16	0e1a954b-851b-493a-8a64-753b15e20a70	cxb4bGmLC_T6RjyKRWMpuDn6W2TyAxvK9LM_q7r2MtaAV85X6RsaxcxPRW1EmuHa	2025-09-13 23:47:41.527947	f	\N			2025-09-06 23:47:41.5283	2025-09-06 23:47:41.5283	\N
fd64d11e-2f74-4afd-8b8e-31a32375c479	0e1a954b-851b-493a-8a64-753b15e20a70	AeDyrAeOg2iEOabRPMs-Pa7gIO5UQjzj2eWo1UHA4SyAVbbw_qdQ-gE85aY_rIyj	2025-09-14 20:49:17.796996	f	\N			2025-09-07 20:49:17.797381	2025-09-07 20:49:17.797381	\N
ca24fe61-1de8-4fee-9bac-a351c17cd7c1	0e1a954b-851b-493a-8a64-753b15e20a70	hffZlDsKlQfq6t_vSlgZA9qi18IctiWXI_sblcpQYQw5dwaB2uxpDKr8-Aj2pY9_	2025-09-13 17:38:19.460141	t	2025-09-07 20:50:41.02175			2025-09-06 17:38:19.460594	2025-09-07 20:50:41.021823	\N
30cbb0a3-f5bd-4513-bd84-0c9d44db9ca4	0e1a954b-851b-493a-8a64-753b15e20a70	mmOhnPzR_Lx2QA32xnNln09M_b7YfdiPV2a0IXvlHzyjNam_b5dUzS9rBVCUjg_q	2025-09-14 20:56:42.102393	t	2025-09-08 21:03:25.147944			2025-09-07 20:56:42.103137	2025-09-08 21:03:25.148442	\N
9a1960a3-31ec-4710-b1cf-0e27905b86bc	0e1a954b-851b-493a-8a64-753b15e20a70	oUFX04bdJLTkxOhAKtXzMET46bhCHsDsCsJ9bP5PzytIUKd2Kl0UuyyjIj1C0383	2025-09-15 21:03:25.164794	f	\N			2025-09-08 21:03:25.165994	2025-09-08 21:03:25.165994	\N
b48168ff-40e2-4c94-bbc6-857886517035	0e1a954b-851b-493a-8a64-753b15e20a70	FuiXuGj8eCIOltP5bj2h1k5QF2dcCzho-k16z8KgVdi9fGE7MVxwT7FD8CvtiVGD	2025-09-14 20:50:41.024648	t	2025-09-08 21:05:12.419155			2025-09-07 20:50:41.024882	2025-09-08 21:05:12.419255	\N
a3902d6f-5195-48f2-bdd9-01b3e3584fd2	0e1a954b-851b-493a-8a64-753b15e20a70	jb8mDkZ2p6qzcxMeQntUID8l5QTr_JoabjtfggwBjB9jvPWNZW_EqodxhJlZbJ4q	2025-09-15 21:05:12.42143	f	\N			2025-09-08 21:05:12.421748	2025-09-08 21:05:12.421748	\N
7e335f9b-f892-4d82-a3f6-289607be9239	0e1a954b-851b-493a-8a64-753b15e20a70	zrGYQ4G3hKisog8SEQa7-KYF5LVHBiNCxbyelZUFFwIyZoDmYxIU3Bg7hs82v2ox	2025-09-15 21:06:42.830231	f	\N			2025-09-08 21:06:42.830697	2025-09-08 21:06:42.830697	\N
00f4cb7c-d8ba-44e4-9f99-c968c1665584	0e1a954b-851b-493a-8a64-753b15e20a70	h0Zgt0VjFzVGfY4CKE5gSBn3cqeacM0ZsDRzRh8QfvUd8ikdDqj5BFJtj4d3vKrj	2025-09-15 21:29:00.787807	f	\N			2025-09-08 21:29:00.789781	2025-09-08 21:29:00.789781	\N
8f446580-1493-4653-9e04-e07d5293217f	0e1a954b-851b-493a-8a64-753b15e20a70	TYsKrs1c05n1nRqLd7xNQZ2lu6LIsr6Q65SDo63KCA7YC1j7ZxmW7JfdGcCJvQfP	2025-09-15 21:29:25.173141	f	\N			2025-09-08 21:29:25.173741	2025-09-08 21:29:25.173741	\N
ee89da24-987c-44de-a5e8-d7aaae42e1ad	0e1a954b-851b-493a-8a64-753b15e20a70	XWDeTQBFrnShsy2cx2XpwQzu9jtVX0JccjYyfwtDotf5VEaFD_rnUc4V2oX_lWnu	2025-09-15 22:01:43.754071	t	2025-09-08 22:06:25.751938			2025-09-08 22:01:43.75798	2025-09-08 22:06:25.752062	\N
d5505a48-3955-4a1b-aeb9-7092fcd1c8a2	0e1a954b-851b-493a-8a64-753b15e20a70	F9DTKZIYvZruRdAnDq25XIrDlL1vrVop06TK7JKqWXUA6i4ukVfHMO6Pbr8YQV-v	2025-09-15 22:06:25.753616	t	2025-09-10 11:44:52.377925			2025-09-08 22:06:25.75382	2025-09-10 11:44:52.378289	\N
5f45a57e-8181-4f64-984f-05be62b7b5f0	0e1a954b-851b-493a-8a64-753b15e20a70	XZnB8KOY7LAW14dILCtaPe4WL1P5G-M5PWKOPUxID2WLBkvcjqrNEu5nSbsTmLSo	2025-09-17 11:44:52.407718	f	\N			2025-09-10 11:44:52.408184	2025-09-10 11:44:52.408184	\N
6ba45dda-0977-4375-939c-e7d2b208b14b	0e1a954b-851b-493a-8a64-753b15e20a70	4VHypacCwhp-38CX48wwFf16XneamI4XNXLlI67o3xhHGrEubM4ZzkhdrFvPWNIm	2025-09-17 12:13:15.282933	f	\N			2025-09-10 12:13:15.283482	2025-09-10 12:13:15.283482	\N
68b2d75f-1cac-4d5c-961d-9026f55ff0cd	0e1a954b-851b-493a-8a64-753b15e20a70	d0I2n4uNSv_3FcTyJtYDICP-cOdUgvG3LG1g9fjBsvUfXZbOE_bMsBZA7JvFdSuQ	2025-09-17 23:10:19.492774	t	2025-09-12 08:20:10.743911			2025-09-10 23:10:19.49441	2025-09-12 08:20:10.744068	\N
8748df51-327f-4e9e-be6d-c9fca5c07d14	0e1a954b-851b-493a-8a64-753b15e20a70	VEzAK9LKzNk_GWVp0rqNrtu3QViaIzQ847knOkjMF38NYc74zZUa8Apm_0RSq2wv	2025-09-19 08:20:10.749811	f	\N			2025-09-12 08:20:10.750081	2025-09-12 08:20:10.750081	\N
0075bbcc-4c09-4e86-aafe-59dfcf213ef5	0e1a954b-851b-493a-8a64-753b15e20a70	6IjSEZgho_8ZsiGfcBoM17BPuZJ4ngXLfvQzDfguyz84DL8ncCURFe4fA8YPZmwC	2025-09-19 09:50:18.801189	f	\N			2025-09-12 09:50:18.801666	2025-09-12 09:50:18.801666	\N
793032b8-40d3-48fb-b4d9-ec8ac1ed018a	0e1a954b-851b-493a-8a64-753b15e20a70	-D3wYiL6WB6Ol67w5_X75u8hp6OlhM1ct6PZwj6go-aC0qBYetVjTiY2Xl7hD-UP	2025-09-19 09:51:43.875758	f	\N			2025-09-12 09:51:43.884903	2025-09-12 09:51:43.884903	\N
51f26d24-e846-4d8e-accb-2036567e89bc	0e1a954b-851b-493a-8a64-753b15e20a70	yznUN64PM9a_WMiDZTdqEb3IBEm90FLRszbMiClK2bj0I7Uiy7boZ6SK1RsK4BcV	2025-09-19 10:32:37.694168	f	\N			2025-09-12 10:32:37.694423	2025-09-12 10:32:37.694423	\N
f63a2efe-547e-4c0b-bf01-195a5da1e6f7	0e1a954b-851b-493a-8a64-753b15e20a70	jU5BEs_cDwAwTdKRhdKRZQ76P9fadHWe4NfxlT2QzziwlM3iLlRPnPZVM9X8oSEb	2025-09-19 10:36:44.888321	f	\N			2025-09-12 10:36:44.88878	2025-09-12 10:36:44.88878	\N
f36b6e36-bbd6-4cf6-97bc-69966764fadc	0e1a954b-851b-493a-8a64-753b15e20a70	wCYPZkLy4ke_A5liMh91r4DlHOgMKFqGc2lVD28Q7zEU68zCi4tMeC449BWnJn4c	2025-09-19 12:17:14.035482	t	2025-09-12 23:24:58.580878			2025-09-12 12:17:14.036088	2025-09-12 23:24:58.580911	\N
d60ab5b3-339f-42b9-971b-ce618b23c749	0e1a954b-851b-493a-8a64-753b15e20a70	Had0Q2IarMjxO9N39sH_Yn-vlOGoxznguPsIdrudQtEPEgWC1rNn8wsfW2Cqhyd6	2025-09-19 23:24:58.584272	f	\N			2025-09-12 23:24:58.585001	2025-09-12 23:24:58.585001	\N
c7257ac9-5622-4c16-8504-a2670ec192b2	0e1a954b-851b-493a-8a64-753b15e20a70	wrBjb5TBcH-KMumHzHhy6NerFrSsfVKnNBwa6OUCgBJ8nPSYWOqkb3kmvXo6bZRa	2025-10-06 15:25:51.552513	f	\N			2025-09-29 15:25:51.553033	2025-09-29 15:25:51.553033	\N
40107e39-9670-4fe6-8e91-b4fa10f7457e	0e1a954b-851b-493a-8a64-753b15e20a70	a_K786WjDADug_39nV6_TR0_6Mo8oOA1nTz_YhB1riQfpRLYo8YxuRKZKVslIYE-	2025-10-07 08:57:30.193834	f	\N			2025-09-30 08:57:30.194052	2025-09-30 08:57:30.194052	\N
73a11f83-4e1d-482f-80ab-952fd0f1dc97	0e1a954b-851b-493a-8a64-753b15e20a70	zxqiMsFzjZQP-U-CZnQOv6j8BzW7HE8rrno6Oyx60KRdji1bzFZ2oBl0l9hLJV4z	2025-10-07 13:01:02.821395	f	\N			2025-09-30 13:01:02.822084	2025-09-30 13:01:02.822084	\N
916385a5-9873-4166-9198-d893d056dbcf	0e1a954b-851b-493a-8a64-753b15e20a70	6-LYUsgdFR5tIZ79l-bX5AgU2NQhwL_6wLZ6dielyVlguyCbzSldKvSK_g1zcHmF	2025-10-07 13:03:30.238391	t	2025-09-30 13:05:42.369205			2025-09-30 13:03:30.238726	2025-09-30 13:05:42.36927	\N
a0048a10-85cf-4be0-b96d-ad5e7d67e4ee	0e1a954b-851b-493a-8a64-753b15e20a70	_AHPXdl9MEr4bdVrIR7b647DvTaEPaWrCSh4Xp0fMcEaItGQ4QXtnxkbdHPqAgFS	2025-10-07 13:05:42.370259	f	\N			2025-09-30 13:05:42.370435	2025-09-30 13:05:42.370435	\N
8df8451e-70ef-459b-aeba-91bd0f49d484	0e1a954b-851b-493a-8a64-753b15e20a70	hefbYEHM9HOe8qGjrpsih3Z3jk_HFj1_orlvmuprVe1cuHPQ8-5OvjPI8k5Z5sm0	2025-10-07 13:14:55.593013	f	\N			2025-09-30 13:14:55.593486	2025-09-30 13:14:55.593486	\N
35a1f26b-1f55-4351-b42a-ebb2cfaae3ef	0e1a954b-851b-493a-8a64-753b15e20a70	B6EP_1aEqeVWrUAaMgrBFZfSAitTixt1SsRbmCgVK4-KvDmHxGxhs6uI5f4MVAUZ	2025-10-16 22:10:59.751187	f	\N			2025-10-09 22:10:59.751587	2025-10-09 22:10:59.751587	\N
24579f18-badc-497b-ad25-5b82cf4b4c49	0e1a954b-851b-493a-8a64-753b15e20a70	MQvVTAh3eH1SAng8p74P2GTG7XUHVL9O5VlG-DtWKSeOeuS3nuhTyZ6RO2HybEnG	2025-10-28 08:36:32.977441	f	\N			2025-10-21 08:36:32.978004	2025-10-21 08:36:32.978004	\N
a5877073-e782-4cf4-9ce6-7ad1207c84ea	0e1a954b-851b-493a-8a64-753b15e20a70	Nw72XW9jkyojt486tTtWjs68uaDNiFp49CMpvIYlPfSxIghjP6eT0kXRX78sx3wv	2025-10-30 22:33:37.587827	f	\N			2025-10-23 22:33:37.588499	2025-10-23 22:33:37.588499	\N
\.


--
-- Data for Name: relationship_analyses; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.relationship_analyses (id, user_id, person_id, connection_strength, engagement_quality, communication_balance, energy_alignment, relationship_health, overall_score, summary, key_insights, patterns, strengths, concerns, trend_direction, provider, model, tokens_used, processing_time_ms, version, analyzed_at, interactions_count, created_at, updated_at, deleted_at) FROM stdin;
59eecb0b-5ab2-4c27-bf9f-b3b7c29f1e48	0e1a954b-851b-493a-8a64-753b15e20a70	9b476068-66aa-44ca-8325-fe80910143ca	0.00	0.00	0.00	0.00	0.00	0.00	Blendina's family relationship is currently facing challenges, with consistently low-quality interactions and varied energy patterns. There is a need for improved communication and more positive engagements to enhance the relationship health.	{}	{"There are frequent neutral and draining interactions, with few energizing moments.","Recent interactions have shown no improvement in quality, which may point to persistent issues."}	{"Blendina is maintaining some level of engagement with family despite challenges."}	{"The consistently low-quality ratings of interactions indicate a need for intervention.","The relationship health score is below average, suggesting potential strain and dissatisfaction."}		openai	gpt-4o	670	5109	1	2025-10-23 23:03:45.939285	6	2025-10-23 23:03:45.940105	2025-10-23 23:03:45.940105	\N
\.


--
-- Data for Name: relationship_statuses; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.relationship_statuses (id, name, color, created_at, updated_at, deleted_at) FROM stdin;
bc29f97e-2542-437c-a01a-adb8a09d7c39	new	blue	2025-08-14 10:43:58.907687	2025-08-14 10:43:58.907687	\N
f29de9cd-18ae-42f8-8849-c33ba654bfd5	stable	green	2025-08-14 10:43:58.907687	2025-08-14 10:43:58.907687	\N
299c5233-e0ff-4c12-83bd-737add383854	fading	orange	2025-08-14 10:43:58.907687	2025-08-14 10:43:58.907687	\N
d330b9db-583f-4a71-b136-7221c64936e9	tense	red	2025-08-14 10:43:58.907687	2025-08-14 10:43:58.907687	\N
\.


--
-- Data for Name: user_consents; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.user_consents (id, user_id, consent_type, granted, version, ip_address, user_agent, granted_at, revoked_at, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: vyve
--

COPY public.users (id, username, email, email_verified, password_hash, avatar_url, display_name, bio, timezone, locale, last_login_at, last_activity_at, streak_count, last_reflection_at, settings, metadata, data_residency, created_at, updated_at, deleted_at, onboarding_completed, onboarding_steps) FROM stdin;
8e80f14b-6072-4930-87a5-0fe86b060750	bob	bob@example.com	t	$2a$10$N9qo8uLOickgx2ZMRZo5i.uG8vH0Hbd.Zo9/6uzRJpcDO/0u5G/	\N	Bob Johnson	\N	UTC	en	\N	\N	0	\N	{}	{}	us	\N	\N	\N	f	[]
6a8d8efc-0111-440a-87b2-4744d0d4bc4d	carol	carol@example.com	t	$2a$10$N9qo8uLOickgx2ZMRZo5i.uG8vH0Hbd.Zo9/6uzRJpcDO/0u5G/	\N	Carol White	\N	UTC	en	\N	\N	0	\N	{}	{}	us	\N	\N	\N	f	[]
0e1a954b-851b-493a-8a64-753b15e20a70	alice	alice@example.com	t	a.uG8vH0Hbd.Zo9/6uzRJpcDO/0u5G/		Luan Jubica	Bio 1			2025-10-23 22:33:37.552787	2025-10-23 22:33:37.552787	0	\N	{"darkMode": true, "notifications": true}	{}	us	2025-08-16 17:35:35.01052	2025-10-23 22:33:37.553238	\N	t	[]
\.


--
-- Name: ai_analysis_jobs ai_analysis_jobs_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.ai_analysis_jobs
    ADD CONSTRAINT ai_analysis_jobs_pkey PRIMARY KEY (id);


--
-- Name: audit_logs audit_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_pkey PRIMARY KEY (id);


--
-- Name: auth_providers auth_providers_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.auth_providers
    ADD CONSTRAINT auth_providers_pkey PRIMARY KEY (id);


--
-- Name: auth_providers auth_providers_provider_provider_id_key; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.auth_providers
    ADD CONSTRAINT auth_providers_provider_provider_id_key UNIQUE (provider, provider_id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: communication_methods communication_methods_name_key; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.communication_methods
    ADD CONSTRAINT communication_methods_name_key UNIQUE (name);


--
-- Name: communication_methods communication_methods_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.communication_methods
    ADD CONSTRAINT communication_methods_pkey PRIMARY KEY (id);


--
-- Name: daily_metrics daily_metrics_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.daily_metrics
    ADD CONSTRAINT daily_metrics_pkey PRIMARY KEY (user_id, date);


--
-- Name: data_exports data_exports_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.data_exports
    ADD CONSTRAINT data_exports_pkey PRIMARY KEY (id);


--
-- Name: energy_patterns energy_patterns_name_key; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.energy_patterns
    ADD CONSTRAINT energy_patterns_name_key UNIQUE (name);


--
-- Name: energy_patterns energy_patterns_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.energy_patterns
    ADD CONSTRAINT energy_patterns_pkey PRIMARY KEY (id);


--
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);


--
-- Name: intentions intentions_name_key; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.intentions
    ADD CONSTRAINT intentions_name_key UNIQUE (name);


--
-- Name: intentions intentions_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.intentions
    ADD CONSTRAINT intentions_pkey PRIMARY KEY (id);


--
-- Name: interactions interactions_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.interactions
    ADD CONSTRAINT interactions_pkey PRIMARY KEY (id);


--
-- Name: nudges nudges_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.nudges
    ADD CONSTRAINT nudges_pkey PRIMARY KEY (id);


--
-- Name: people people_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.people
    ADD CONSTRAINT people_pkey PRIMARY KEY (id);


--
-- Name: push_tokens push_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.push_tokens
    ADD CONSTRAINT push_tokens_pkey PRIMARY KEY (id);


--
-- Name: push_tokens push_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.push_tokens
    ADD CONSTRAINT push_tokens_token_key UNIQUE (token);


--
-- Name: reflections reflections_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.reflections
    ADD CONSTRAINT reflections_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_token_key UNIQUE (token);


--
-- Name: relationship_analyses relationship_analyses_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.relationship_analyses
    ADD CONSTRAINT relationship_analyses_pkey PRIMARY KEY (id);


--
-- Name: relationship_statuses relationship_statuses_name_key; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.relationship_statuses
    ADD CONSTRAINT relationship_statuses_name_key UNIQUE (name);


--
-- Name: relationship_statuses relationship_statuses_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.relationship_statuses
    ADD CONSTRAINT relationship_statuses_pkey PRIMARY KEY (id);


--
-- Name: categories uq_categories_user_name; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT uq_categories_user_name UNIQUE (user_id, name);


--
-- Name: user_consents user_consents_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.user_consents
    ADD CONSTRAINT user_consents_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: idx_ai_analysis_jobs_created_at; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_ai_analysis_jobs_created_at ON public.ai_analysis_jobs USING btree (created_at DESC);


--
-- Name: idx_ai_analysis_jobs_deleted_at; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_ai_analysis_jobs_deleted_at ON public.ai_analysis_jobs USING btree (deleted_at);


--
-- Name: idx_ai_analysis_jobs_priority; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_ai_analysis_jobs_priority ON public.ai_analysis_jobs USING btree (priority DESC);


--
-- Name: idx_ai_analysis_jobs_status; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_ai_analysis_jobs_status ON public.ai_analysis_jobs USING btree (status);


--
-- Name: idx_ai_analysis_jobs_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_ai_analysis_jobs_user_id ON public.ai_analysis_jobs USING btree (user_id);


--
-- Name: idx_audit_logs_action; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_audit_logs_action ON public.audit_logs USING btree (action);


--
-- Name: idx_audit_logs_created_at; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_audit_logs_created_at ON public.audit_logs USING btree (created_at);


--
-- Name: idx_audit_logs_entity; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_audit_logs_entity ON public.audit_logs USING btree (entity_type, entity_id);


--
-- Name: idx_audit_logs_request_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_audit_logs_request_id ON public.audit_logs USING btree (request_id);


--
-- Name: idx_audit_logs_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_audit_logs_user_id ON public.audit_logs USING btree (user_id) WHERE (deleted_at IS NULL);


--
-- Name: idx_auth_providers_provider; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_auth_providers_provider ON public.auth_providers USING btree (provider, provider_id);


--
-- Name: idx_auth_providers_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_auth_providers_user_id ON public.auth_providers USING btree (user_id);


--
-- Name: idx_daily_metrics_date; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_daily_metrics_date ON public.daily_metrics USING btree (date);


--
-- Name: idx_daily_metrics_user_date; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_daily_metrics_user_date ON public.daily_metrics USING btree (user_id, date DESC);


--
-- Name: idx_daily_metrics_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_daily_metrics_user_id ON public.daily_metrics USING btree (user_id);


--
-- Name: idx_data_exports_expires_at; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_data_exports_expires_at ON public.data_exports USING btree (expires_at);


--
-- Name: idx_data_exports_status; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_data_exports_status ON public.data_exports USING btree (status);


--
-- Name: idx_data_exports_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_data_exports_user_id ON public.data_exports USING btree (user_id) WHERE (deleted_at IS NULL);


--
-- Name: idx_events_created_at; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_events_created_at ON public.events USING btree (created_at);


--
-- Name: idx_events_event_type; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_events_event_type ON public.events USING btree (event_type);


--
-- Name: idx_events_session_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_events_session_id ON public.events USING btree (session_id);


--
-- Name: idx_events_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_events_user_id ON public.events USING btree (user_id) WHERE (deleted_at IS NULL);


--
-- Name: idx_events_user_type_date; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_events_user_type_date ON public.events USING btree (user_id, event_type, created_at);


--
-- Name: idx_interactions_energy_impact; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_interactions_energy_impact ON public.interactions USING btree (energy_impact);


--
-- Name: idx_interactions_interaction_at; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_interactions_interaction_at ON public.interactions USING btree (interaction_at);


--
-- Name: idx_interactions_person_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_interactions_person_id ON public.interactions USING btree (person_id) WHERE (deleted_at IS NULL);


--
-- Name: idx_interactions_user_date; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_interactions_user_date ON public.interactions USING btree (user_id, interaction_at);


--
-- Name: idx_interactions_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_interactions_user_id ON public.interactions USING btree (user_id) WHERE (deleted_at IS NULL);


--
-- Name: idx_nudges_acted_on; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_nudges_acted_on ON public.nudges USING btree (acted_on);


--
-- Name: idx_nudges_analysis_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_nudges_analysis_id ON public.nudges USING btree (analysis_id);


--
-- Name: idx_nudges_expires_at; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_nudges_expires_at ON public.nudges USING btree (expires_at);


--
-- Name: idx_nudges_person_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_nudges_person_id ON public.nudges USING btree (person_id) WHERE (deleted_at IS NULL);


--
-- Name: idx_nudges_priority; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_nudges_priority ON public.nudges USING btree (priority);


--
-- Name: idx_nudges_seen; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_nudges_seen ON public.nudges USING btree (seen);


--
-- Name: idx_nudges_source; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_nudges_source ON public.nudges USING btree (source);


--
-- Name: idx_nudges_status; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_nudges_status ON public.nudges USING btree (status);


--
-- Name: idx_nudges_timing; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_nudges_timing ON public.nudges USING btree (timing);


--
-- Name: idx_nudges_type; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_nudges_type ON public.nudges USING btree (type);


--
-- Name: idx_nudges_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_nudges_user_id ON public.nudges USING btree (user_id) WHERE (deleted_at IS NULL);


--
-- Name: idx_people_category_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_people_category_id ON public.people USING btree (category_id);


--
-- Name: idx_people_comm_method_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_people_comm_method_id ON public.people USING btree (communication_method_id);


--
-- Name: idx_people_energy_pattern_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_people_energy_pattern_id ON public.people USING btree (energy_pattern_id);


--
-- Name: idx_people_health_score; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_people_health_score ON public.people USING btree (health_score);


--
-- Name: idx_people_last_interaction; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_people_last_interaction ON public.people USING btree (last_interaction_at);


--
-- Name: idx_people_next_reminder; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_people_next_reminder ON public.people USING btree (next_reminder_at);


--
-- Name: idx_people_relationship_status_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_people_relationship_status_id ON public.people USING btree (relationship_status_id);


--
-- Name: idx_people_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_people_user_id ON public.people USING btree (user_id) WHERE (deleted_at IS NULL);


--
-- Name: idx_push_tokens_active; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_push_tokens_active ON public.push_tokens USING btree (active);


--
-- Name: idx_push_tokens_token; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_push_tokens_token ON public.push_tokens USING btree (token);


--
-- Name: idx_push_tokens_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_push_tokens_user_id ON public.push_tokens USING btree (user_id);


--
-- Name: idx_reflections_completed_at; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_reflections_completed_at ON public.reflections USING btree (completed_at);


--
-- Name: idx_reflections_mood; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_reflections_mood ON public.reflections USING btree (mood);


--
-- Name: idx_reflections_user_date; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_reflections_user_date ON public.reflections USING btree (user_id, completed_at);


--
-- Name: idx_reflections_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_reflections_user_id ON public.reflections USING btree (user_id) WHERE (deleted_at IS NULL);


--
-- Name: idx_refresh_tokens_expires_at; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_refresh_tokens_expires_at ON public.refresh_tokens USING btree (expires_at);


--
-- Name: idx_refresh_tokens_token; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_refresh_tokens_token ON public.refresh_tokens USING btree (token);


--
-- Name: idx_refresh_tokens_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_refresh_tokens_user_id ON public.refresh_tokens USING btree (user_id);


--
-- Name: idx_relationship_analyses_analyzed_at; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_relationship_analyses_analyzed_at ON public.relationship_analyses USING btree (analyzed_at DESC);


--
-- Name: idx_relationship_analyses_deleted_at; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_relationship_analyses_deleted_at ON public.relationship_analyses USING btree (deleted_at);


--
-- Name: idx_relationship_analyses_person_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_relationship_analyses_person_id ON public.relationship_analyses USING btree (person_id);


--
-- Name: idx_relationship_analyses_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_relationship_analyses_user_id ON public.relationship_analyses USING btree (user_id);


--
-- Name: idx_relationship_analyses_user_person; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_relationship_analyses_user_person ON public.relationship_analyses USING btree (user_id, person_id);


--
-- Name: idx_user_consents_granted; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_user_consents_granted ON public.user_consents USING btree (granted);


--
-- Name: idx_user_consents_type; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_user_consents_type ON public.user_consents USING btree (consent_type);


--
-- Name: idx_user_consents_user_id; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_user_consents_user_id ON public.user_consents USING btree (user_id) WHERE (deleted_at IS NULL);


--
-- Name: idx_users_data_residency; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_users_data_residency ON public.users USING btree (data_residency);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_users_email ON public.users USING btree (email) WHERE (deleted_at IS NULL);


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: vyve
--

CREATE INDEX idx_users_username ON public.users USING btree (username) WHERE (deleted_at IS NULL);


--
-- Name: ai_analysis_jobs update_ai_analysis_jobs_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_ai_analysis_jobs_updated_at BEFORE UPDATE ON public.ai_analysis_jobs FOR EACH ROW EXECUTE FUNCTION public.update_ai_tables_updated_at();


--
-- Name: audit_logs update_audit_logs_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_audit_logs_updated_at BEFORE UPDATE ON public.audit_logs FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: auth_providers update_auth_providers_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_auth_providers_updated_at BEFORE UPDATE ON public.auth_providers FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: daily_metrics update_daily_metrics_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_daily_metrics_updated_at BEFORE UPDATE ON public.daily_metrics FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: data_exports update_data_exports_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_data_exports_updated_at BEFORE UPDATE ON public.data_exports FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: events update_events_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON public.events FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: interactions update_interactions_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_interactions_updated_at BEFORE UPDATE ON public.interactions FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: nudges update_nudges_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_nudges_updated_at BEFORE UPDATE ON public.nudges FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: people update_people_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_people_updated_at BEFORE UPDATE ON public.people FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: interactions update_person_last_interaction_trigger; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_person_last_interaction_trigger AFTER INSERT ON public.interactions FOR EACH ROW EXECUTE FUNCTION public.update_person_last_interaction();


--
-- Name: push_tokens update_push_tokens_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_push_tokens_updated_at BEFORE UPDATE ON public.push_tokens FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: reflections update_reflections_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_reflections_updated_at BEFORE UPDATE ON public.reflections FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: refresh_tokens update_refresh_tokens_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_refresh_tokens_updated_at BEFORE UPDATE ON public.refresh_tokens FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: user_consents update_user_consents_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_user_consents_updated_at BEFORE UPDATE ON public.user_consents FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: users update_users_updated_at; Type: TRIGGER; Schema: public; Owner: vyve
--

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: ai_analysis_jobs ai_analysis_jobs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.ai_analysis_jobs
    ADD CONSTRAINT ai_analysis_jobs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: audit_logs audit_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: auth_providers auth_providers_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.auth_providers
    ADD CONSTRAINT auth_providers_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: categories categories_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: daily_metrics daily_metrics_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.daily_metrics
    ADD CONSTRAINT daily_metrics_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: data_exports data_exports_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.data_exports
    ADD CONSTRAINT data_exports_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: events events_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: interactions interactions_person_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.interactions
    ADD CONSTRAINT interactions_person_id_fkey FOREIGN KEY (person_id) REFERENCES public.people(id) ON DELETE CASCADE;


--
-- Name: interactions interactions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.interactions
    ADD CONSTRAINT interactions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: nudges nudges_analysis_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.nudges
    ADD CONSTRAINT nudges_analysis_id_fkey FOREIGN KEY (analysis_id) REFERENCES public.relationship_analyses(id) ON DELETE SET NULL;


--
-- Name: nudges nudges_person_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.nudges
    ADD CONSTRAINT nudges_person_id_fkey FOREIGN KEY (person_id) REFERENCES public.people(id) ON DELETE CASCADE;


--
-- Name: nudges nudges_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.nudges
    ADD CONSTRAINT nudges_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: people people_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.people
    ADD CONSTRAINT people_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE SET NULL;


--
-- Name: people people_communication_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.people
    ADD CONSTRAINT people_communication_method_id_fkey FOREIGN KEY (communication_method_id) REFERENCES public.communication_methods(id) ON DELETE SET NULL;


--
-- Name: people people_energy_pattern_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.people
    ADD CONSTRAINT people_energy_pattern_id_fkey FOREIGN KEY (energy_pattern_id) REFERENCES public.energy_patterns(id) ON DELETE SET NULL;


--
-- Name: people people_intention_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.people
    ADD CONSTRAINT people_intention_id_fkey FOREIGN KEY (intention_id) REFERENCES public.intentions(id) ON DELETE SET NULL;


--
-- Name: people people_relationship_status_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.people
    ADD CONSTRAINT people_relationship_status_id_fkey FOREIGN KEY (relationship_status_id) REFERENCES public.relationship_statuses(id) ON DELETE SET NULL;


--
-- Name: people people_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.people
    ADD CONSTRAINT people_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: push_tokens push_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.push_tokens
    ADD CONSTRAINT push_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: reflections reflections_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.reflections
    ADD CONSTRAINT reflections_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: refresh_tokens refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: relationship_analyses relationship_analyses_person_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.relationship_analyses
    ADD CONSTRAINT relationship_analyses_person_id_fkey FOREIGN KEY (person_id) REFERENCES public.people(id) ON DELETE CASCADE;


--
-- Name: relationship_analyses relationship_analyses_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.relationship_analyses
    ADD CONSTRAINT relationship_analyses_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: user_consents user_consents_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: vyve
--

ALTER TABLE ONLY public.user_consents
    ADD CONSTRAINT user_consents_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict 7LGbFidJ9dXINV6qGK8WtxaikiEKfaIeemyhRVZIRtRconJiZjgAIhHVKud6sCe

