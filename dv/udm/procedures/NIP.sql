/*
=============================================
Author:				Michael Kaiser
Create date:		2025-12-01
Description:		Table representing the NIP-Dimension
Modification:
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.nip_load(p_force boolean DEFAULT FALSE)
    LANGUAGE plpgsql
AS
$proc$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'udm') THEN
        RAISE EXCEPTION $$Schema udm does not exist$$;
    END IF;

    IF p_force THEN
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'nip') THEN
            EXECUTE $$DROP TABLE udm.Nip$$;
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'nip') THEN
        EXECUTE $$CREATE TABLE udm.Nip AS
                      SELECT 'neo4j'              AS record_src,
                             _id                  AS nip_id,
                             name                 AS nip_name,
                             CONCAT('NIP-', name) AS nip_full_name
                      FROM rdv.dump_neo4j_sat_cur
                      WHERE _labels = ':NIP'$$;
        EXECUTE $$CREATE INDEX IF NOT EXISTS idx_udm_nip_nip_id ON udm.Nip (nip_id)$$;
    ELSE
        BEGIN
            EXECUTE $$TRUNCATE TABLE udm.Nip$$;
            EXECUTE $$INSERT INTO udm.Nip (record_src, nip_id, nip_name, nip_full_name)
                          SELECT 'neo4j'              AS record_src,
                                 _id                  AS nip_id,
                                 name                 AS nip_name,
                                 CONCAT('NIP-', name) AS nip_full_name
                          FROM rdv.dump_neo4j_sat_cur
                          WHERE _labels = ':NIP'$$;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE $$Error while refreshing udm.Nip: %$$, SQLERRM;
                RAISE;
        END;
    END IF;

END;
$proc$;