/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Relationship table for relay to software usage
Modification:
2026-01-04  mkaiser     Fix the mix of software and owner
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.relaysoftware_load(p_force boolean DEFAULT FALSE)
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
                     AND table_name = 'relaysoftware') THEN
            EXECUTE 'DROP TABLE udm.RelaySoftware';
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1
                   FROM information_schema.tables
                   WHERE table_schema = 'udm'
                     AND table_name = 'relaysoftware') THEN
        EXECUTE 'CREATE TABLE udm.RelaySoftware AS SELECT ''neo4j'' AS record_src, _start AS relay_id, _end AS software_id FROM rdv.dump_neo4j_sat_cur WHERE _type = ''USES_SOFTWARE''';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relaysoftware_relay_id ON udm.RelaySoftware (relay_id)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relaysoftware_software_id ON udm.RelaySoftware (software_id)';
    ELSE
        BEGIN
            EXECUTE 'TRUNCATE TABLE udm.RelaySoftware';
            EXECUTE 'INSERT INTO udm.RelaySoftware (record_src, relay_id, software_id) SELECT ''neo4j'' AS record_src, _start AS relay_id, _end AS software_id FROM rdv.dump_neo4j_sat_cur WHERE _type = ''USES_SOFTWARE''';
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Error while refreshing udm.RelaySoftware: %', SQLERRM;
                RAISE;
        END;
    END IF;

END;
$proc$;
