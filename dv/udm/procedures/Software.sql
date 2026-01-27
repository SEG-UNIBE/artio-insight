/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Table representing the software dimension
Modification:
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.software_load(p_force boolean DEFAULT FALSE)
    LANGUAGE plpgsql
AS
$proc$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'udm') THEN
        RAISE EXCEPTION 'Schema udm does not exist';
    END IF;

    IF p_force THEN
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'software') THEN
            EXECUTE 'DROP TABLE udm.Software';
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'software') THEN
        EXECUTE $$CREATE TABLE udm.Software AS
                      SELECT 'neo4j'  AS record_src,
                             _id      AS software_id,
                             software AS software
                      FROM rdv.dump_neo4j_sat_cur
                      WHERE _labels = ':Software'$$;
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_software_software_id ON udm.Software (software_id)';
    ELSE
        BEGIN
            EXECUTE 'TRUNCATE TABLE udm.Software';
            EXECUTE $$INSERT INTO udm.Software (record_src, software_id, software)
                          SELECT 'neo4j'  AS record_src,
                                 _id      AS software_id,
                                 software AS software
                          FROM rdv.dump_neo4j_sat_cur
                          WHERE _labels = ':Software'$$;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Error while refreshing udm.Software: %', SQLERRM;
                RAISE;
        END;
    END IF;

END;
$proc$;
