/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Relationship table for relay to relayaltname usage
Modification:
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.relayrelayaltname_load(p_force boolean DEFAULT FALSE)
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
                     AND table_name = 'relayrelayaltname') THEN
            EXECUTE 'DROP TABLE udm.RelayRelayAltName';
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1
                   FROM information_schema.tables
                   WHERE table_schema = 'udm'
                     AND table_name = 'relayrelayaltname') THEN
        EXECUTE 'CREATE TABLE udm.RelayRelayAltName AS SELECT ''neo4j'' AS record_src, _start AS relay_id, _end AS relayaltname_id FROM rdv.dump_neo4j_sat_cur WHERE _type = ''ALT_NAME''';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relayrelayaltname_relay_id ON udm.RelayRelayAltName (relay_id)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relayrelayaltname_relayaltname_id ON udm.RelayRelayAltName (relayaltname_id)';
    ELSE
        BEGIN
            EXECUTE 'TRUNCATE TABLE udm.RelayRelayAltName';
            EXECUTE 'INSERT INTO udm.RelayRelayAltName (record_src, relay_id, relayaltname_id) SELECT ''neo4j'' AS record_src, _start AS relay_id, _end AS relayaltname_id FROM rdv.dump_neo4j_sat_cur WHERE _type = ''ALT_NAME''';
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Error while refreshing udm.RelayRelayAltName: %', SQLERRM;
                RAISE;
        END;
    END IF;

END;
$proc$;
