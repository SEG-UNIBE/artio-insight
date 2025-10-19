/*
=============================================
Author:				Michael Kaiser
Create date:		2025-08-20
Description:		Script to setup the tables in the inb layer
Modification:
2025-09-22	mkaiser     Adding tables for artioinf
2025-10-19  mkaiser     Adding SET UNLOGGED statements for inb tables
=============================================
*/

DROP TABLE IF EXISTS inb.events_artiostr;

CREATE TABLE IF NOT EXISTS inb.events_artiostr
(
    id         bigint                            NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    created    bigint                            NOT NULL,
    event_id   text COLLATE pg_catalog."default",
    pubkey     text COLLATE pg_catalog."default" NOT NULL,
    kind       bigint                            NOT NULL,
    content    text COLLATE pg_catalog."default" NOT NULL,
    sig        text COLLATE pg_catalog."default" NOT NULL,
    tags       jsonb                             NOT NULL DEFAULT '[]'::jsonb,
    CONSTRAINT events_artiostr_pkey PRIMARY KEY (id)
);

ALTER TABLE inb.events_artiostr SET UNLOGGED;

-- Table: public.logs
DROP TABLE IF EXISTS inb.logs_artiostr;

CREATE TABLE IF NOT EXISTS inb.logs_artiostr
(
    id         bigint                            NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    ip         text COLLATE pg_catalog."default",
    type       text COLLATE pg_catalog."default" NOT NULL,
    content    text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT logs_artiostr_pkey PRIMARY KEY (id)
);

ALTER TABLE inb.logs_artiostr SET UNLOGGED;

DROP TABLE IF EXISTS inb.events_artioinf;

CREATE TABLE IF NOT EXISTS inb.events_artioinf
(
    id         bigint                            NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    created    bigint                            NOT NULL,
    event_id   text COLLATE pg_catalog."default",
    pubkey     text COLLATE pg_catalog."default" NOT NULL,
    kind       bigint                            NOT NULL,
    content    text COLLATE pg_catalog."default" NOT NULL,
    sig        text COLLATE pg_catalog."default" NOT NULL,
    tags       jsonb                             NOT NULL DEFAULT '[]'::jsonb,
    CONSTRAINT events_artioinf_pkey PRIMARY KEY (id)
);

ALTER TABLE inb.events_artioinf SET UNLOGGED;

-- Table: public.logs
DROP TABLE IF EXISTS inb.logs_artioinf;

CREATE TABLE IF NOT EXISTS inb.logs_artioinf
(
    id         bigint                            NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    ip         text COLLATE pg_catalog."default",
    type       text COLLATE pg_catalog."default" NOT NULL,
    content    text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT logs_artioinf_pkey PRIMARY KEY (id)
);

ALTER TABLE inb.logs_artioinf SET UNLOGGED;