/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Table representing the RelayAlternativeName dimension
Modification:
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/


CREATE OR REPLACE PROCEDURE udm.relayalternatename_load(p_force boolean DEFAULT FALSE)
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
                     AND table_name = 'relayalternativename') THEN
            EXECUTE 'DROP TABLE udm.RelayAlternativeName';
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1
                   FROM information_schema.tables
                   WHERE table_schema = 'udm'
                     AND table_name = 'relayalternativename') THEN
        EXECUTE $$CREATE TABLE udm.RelayAlternativeName AS
                      SELECT 'neo4j' AS record_src,
                             _id     AS relayaltname_id,
                             name    AS relayaltname
                      FROM rdv.dump_neo4j_sat_cur
                      WHERE _labels = ':RelayAlternativeName'$$;
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relayaltname_relayaltname_id ON udm.RelayAlternativeName (relayaltname_id)';
    ELSE
        BEGIN
            EXECUTE 'TRUNCATE TABLE udm.RelayAlternativeName';
            EXECUTE $$INSERT INTO udm.RelayAlternativeName (record_src, relayaltname_id, relayaltname)
                          SELECT 'neo4j' AS record_src,
                                 _id     AS relayaltname_id,
                                 name    AS relayaltname
                          FROM rdv.dump_neo4j_sat_cur
                          WHERE _labels = ':RelayAlternativeName'$$;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Error while refreshing udm.RelayAlternativeName: %', SQLERRM;
                RAISE;
        END;
    END IF;

END;
$proc$;
