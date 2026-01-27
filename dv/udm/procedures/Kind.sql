/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Table for the kinds of events
Modification:
2025-11-25  mkaiser     Changing the script to physical table
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.kind_load(p_force boolean DEFAULT FALSE)
    LANGUAGE plpgsql
AS
$proc$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'udm') THEN
        RAISE EXCEPTION 'Schema udm does not exist';
    END IF;

    IF p_force THEN
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'kind') THEN
            EXECUTE 'DROP TABLE udm.kind';
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'kind') THEN
        EXECUTE $$CREATE TABLE udm.kind AS
                      SELECT record_src,
                             kind,
                             COUNT(DISTINCT pubkey) AS user_count,
                             COUNT(event_id)        AS event_count
                      FROM udm.events_aggregation
                      GROUP BY kind, record_src$$;
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_kind_kind ON udm.kind (kind)';
    ELSE
        BEGIN
            EXECUTE 'TRUNCATE TABLE udm.kind';
            EXECUTE $$INSERT INTO udm.kind (record_src, kind, user_count, event_count)
                          SELECT record_src,
                                 kind,
                                 COUNT(DISTINCT pubkey) AS user_count,
                                 COUNT(event_id)        AS event_count
                          FROM udm.events_aggregation
                          GROUP BY kind, record_src$$;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Error while refreshing udm.kind: %', SQLERRM;
                RAISE;
        END;
    END IF;
END;
$proc$;
