/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Relationship table for relay to relay detection usage
Modification:
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.relaydetection_load(p_force boolean DEFAULT FALSE)
    LANGUAGE plpgsql
AS
$proc$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'udm') THEN
        RAISE EXCEPTION 'Schema udm does not exist';
    END IF;

    IF p_force THEN
        IF EXISTS (SELECT 1
                   FROM information_schema.tables
                   WHERE table_schema = 'udm'
                     AND table_name = 'relaydetection') THEN
            EXECUTE 'DROP TABLE udm.RelayDetection';
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1
                   FROM information_schema.tables
                   WHERE table_schema = 'udm'
                     AND table_name = 'relaydetection') THEN
        EXECUTE $$CREATE TABLE udm.RelayDetection AS
                      SELECT 'neo4j' AS record_src,
                             _start  AS relay_id,
                             _end    AS detected_relay_id
                      FROM rdv.dump_neo4j_sat_cur
                      WHERE _type = 'DETECTED'$$;
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relaydetection_relay_id ON udm.RelayDetection (relay_id)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relaydetection_relay_detected_id ON udm.RelayDetection (detected_relay_id)';
    ELSE
        BEGIN
            EXECUTE 'TRUNCATE TABLE udm.RelayDetection';
            EXECUTE $$INSERT INTO udm.RelayDetection (record_src, relay_id, detected_relay_id)
                          SELECT 'neo4j' AS record_src,
                                 _start  AS relay_id,
                                 _end    AS detected_relay_id
                          FROM rdv.dump_neo4j_sat_cur
                          WHERE _type = 'DETECTED'$$;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Error while refreshing udm.RelayDetection: %', SQLERRM;
                RAISE;
        END;
    END IF;

END;
$proc$;
