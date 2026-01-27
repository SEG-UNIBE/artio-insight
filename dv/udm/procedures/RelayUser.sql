/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Relationship table for relay to user usage
Modification:
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.relayuser_load(p_force boolean DEFAULT FALSE)
    LANGUAGE plpgsql
AS
$proc$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'udm') THEN
        RAISE EXCEPTION 'Schema udm does not exist';
    END IF;

    IF p_force THEN
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'relayuser') THEN
            EXECUTE 'DROP TABLE udm.RelayUser';
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'relayuser') THEN
        EXECUTE 'CREATE TABLE udm.RelayUser AS SELECT ''neo4j'' AS record_src, _start AS user_id, _end AS relay_id FROM rdv.dump_neo4j_sat_cur WHERE _type = ''USES''';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relayuser_relay_id ON udm.RelayUser (relay_id)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relayuser_user_id ON udm.RelayUser (user_id)';
    ELSE
        BEGIN
            EXECUTE 'TRUNCATE TABLE udm.RelayUser';
            EXECUTE 'INSERT INTO udm.RelayUser (record_src, user_id, relay_id) SELECT ''neo4j'' AS record_src, _start AS user_id, _end AS relay_id FROM rdv.dump_neo4j_sat_cur WHERE _type = ''USES''';
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Error while refreshing udm.RelayUser: %', SQLERRM;
                RAISE;
        END;
    END IF;

END;
$proc$;
