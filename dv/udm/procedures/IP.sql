/*
=============================================
Author:				Michael Kaiser
Create date:		2026-01-09
Description:		Table representing the IP dimension
Modification:
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.ip_load(p_force boolean DEFAULT FALSE)
    LANGUAGE plpgsql
AS
$proc$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'udm') THEN
        RAISE EXCEPTION $$Schema udm does not exist$$;
    END IF;

    IF p_force THEN
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'ip') THEN
            EXECUTE $$DROP TABLE udm.IP$$;
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'ip') THEN
        EXECUTE $$CREATE TABLE udm.IP AS
                      SELECT 'neo4j'  AS record_src,
                             _id      AS ip_id,
                             address AS ip_address,
                             CAST(address AS inet) as ip,
                             family(CAST(address AS inet)) as ip_type
                      FROM rdv.dump_neo4j_sat_cur
                      WHERE _labels = ':IP'$$;
        EXECUTE $$CREATE INDEX IF NOT EXISTS idx_udm_ip_ip_id ON udm.IP (ip_id)$$;
        EXECUTE $$CREATE INDEX IF NOT EXISTS idx_udm_ip_ip_address ON udm.IP (ip_address)$$;
    ELSE
        BEGIN
            EXECUTE $$TRUNCATE TABLE udm.IP$$;
            EXECUTE $$INSERT INTO udm.IP (record_src, ip_id, ip_address, ip, ip_type)
                          SELECT 'neo4j'  AS record_src,
                                 _id      AS ip_id,
                                 address AS ip_address,
                                 CAST(address AS inet) as ip,
                                 family(CAST(address AS inet)) as ip_type
                          FROM rdv.dump_neo4j_sat_cur
                          WHERE _labels = ':IP'$$;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE $$Error while refreshing udm.IP: %$$, SQLERRM;
                RAISE;
        END;
    END IF;
END;
$proc$;
