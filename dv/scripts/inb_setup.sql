/*
=============================================
Author:				Michael Kaiser
Create date:		2025-08-20
Description:		Script to setup the tables in the inb layer
Modification:
YYYY-MM-DD	mkaiser
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
    CONSTRAINT events_pkey PRIMARY KEY (id)
);

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
    CONSTRAINT logs_pkey PRIMARY KEY (id)
);