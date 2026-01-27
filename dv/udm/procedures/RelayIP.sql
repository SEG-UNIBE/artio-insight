/*
=============================================
Author:				Michael Kaiser
Create date:		2026-01-09
Description:		Relationship table for relay to relay ip relationship
Modification:
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.relayip_load(p_force boolean DEFAULT FALSE)
    LANGUAGE plpgsql
AS
$proc$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'udm') THEN
        RAISE EXCEPTION 'Schema udm does not exist';
    END IF;

    IF p_force THEN
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'relayip') THEN
            EXECUTE 'DROP TABLE udm.RelayIp';
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'relayip') THEN
        EXECUTE $$CREATE TABLE udm.RelayIp AS
                      SELECT 'neo4j' AS record_src,
                             _start  AS relay_id,
                             _end    AS ip_id
                      FROM rdv.dump_neo4j_sat_cur
                      WHERE _type = 'HAS_IP'$$;
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relayip_relay_id ON udm.RelayIp (relay_id)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relayip_ip_id ON udm.RelayIp (ip_id)';
    ELSE
        BEGIN
            EXECUTE 'TRUNCATE TABLE udm.RelayIp';
            EXECUTE $$INSERT INTO udm.RelayIp (record_src, relay_id, ip_id)
                          SELECT 'neo4j' AS record_src,
                                 _start  AS relay_id,
                                 _end    AS ip_id
                          FROM rdv.dump_neo4j_sat_cur
                          WHERE _type = 'HAS_IP'$$;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Error while refreshing udm.RelayIp: %', SQLERRM;
                RAISE;
        END;
    END IF;

END;
$proc$;
