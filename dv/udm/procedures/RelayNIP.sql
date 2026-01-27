/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Relationship table for relay to relay nip relationship
Modification:
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.relaynip_load(p_force boolean DEFAULT FALSE)
    LANGUAGE plpgsql
AS
$proc$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'udm') THEN
        RAISE EXCEPTION 'Schema udm does not exist';
    END IF;

    IF p_force THEN
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'relaynip') THEN
            EXECUTE 'DROP TABLE udm.RelayNip';
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'relaynip') THEN
        EXECUTE $$CREATE TABLE udm.RelayNip AS
                      SELECT 'neo4j' AS record_src,
                             _start  AS relay_id,
                             _end    AS nip_id
                      FROM rdv.dump_neo4j_sat_cur
                      WHERE _type = 'IMPLEMENTS'$$;
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relaynip_relay_id ON udm.RelayNip (relay_id)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relaynip_nip_id ON udm.RelayNip (nip_id)';
    ELSE
        BEGIN
            EXECUTE 'TRUNCATE TABLE udm.RelayNip';
            EXECUTE $$INSERT INTO udm.RelayNip (record_src, relay_id, nip_id)
                          SELECT 'neo4j' AS record_src,
                                 _start  AS relay_id,
                                 _end    AS nip_id
                          FROM rdv.dump_neo4j_sat_cur
                          WHERE _type = 'IMPLEMENTS'$$;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Error while refreshing udm.RelayNip: %', SQLERRM;
                RAISE;
        END;
    END IF;

END;
$proc$;
