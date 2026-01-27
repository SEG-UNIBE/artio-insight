/*
=============================================
Author:				Michael Kaiser
Create date:		2025-12-01
Description:		Relationship table for relay to nip kinds relationship
Modification:
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.nipkind_load(p_force boolean DEFAULT FALSE)
    LANGUAGE plpgsql
AS
$proc$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'udm') THEN
        RAISE EXCEPTION 'Schema udm does not exist';
    END IF;

    IF p_force THEN
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'nipkind') THEN
            EXECUTE 'DROP TABLE udm.NipKind';
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'nipkind') THEN
        EXECUTE $$CREATE TABLE udm.NipKind AS
                      SELECT 'nostrbook.dev' AS record_src,
                             nip_name  AS nip_name,
                             kind    AS kind
                      FROM rdv.nipkinds_nipfile_sat_cur$$;
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_nipkind_nip_name ON udm.NipKind (nip_name)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_nipkind_kind ON udm.NipKind (kind)';

    ELSE
        BEGIN
            EXECUTE 'TRUNCATE TABLE udm.NipKind';
            EXECUTE $$INSERT INTO udm.NipKind (record_src, nip_name, kind)
                          SELECT 'nostrbook.dev' AS record_src,
                                 nip_name  AS nip_name,
                                 kind    AS kind
                          FROM rdv.nipkinds_nipfile_sat_cur$$;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Error while refreshing udm.NipKind: %', SQLERRM;
                RAISE;
        END;
    END IF;
END;
$proc$;
