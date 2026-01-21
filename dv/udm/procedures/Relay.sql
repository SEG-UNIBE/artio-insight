/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Table for getting stats about the relays
Modification:
2025-11-25  mkaiser     Changing the script to physical table and including neo4j data
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.relay_load(p_force boolean DEFAULT FALSE)
    LANGUAGE plpgsql
AS
$proc$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'udm') THEN
        RAISE EXCEPTION 'Schema udm does not exist';
    END IF;

    IF p_force THEN
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'relay') THEN
            EXECUTE 'DROP TABLE udm.Relay';
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'relay') THEN
        EXECUTE $$CREATE TABLE udm.Relay AS
                      WITH relay_prep AS (SELECT DISTINCT record_src
                                          FROM udm.events_aggregation
                                          UNION
                                          SELECT DISTINCT record_src
                                          FROM udm.logs_aggregation)
                      SELECT relay_prep.record_src,
                             NULL                             AS relay_id,
                             relay_prep.record_src            AS relay_name,
                             COUNT(event_agg.event_id)        AS event_count,
                             COALESCE(MAX(deleted_ev.cnt), 0) AS event_deleted_count,
                             COUNT(DISTINCT event_agg.pubkey) AS pubkey_count
                      FROM relay_prep
                               LEFT JOIN udm.events_aggregation event_agg ON event_agg.record_src = relay_prep.record_src
                               LEFT JOIN (SELECT record_src, COUNT(event_id) AS cnt
                                          FROM udm.events_aggregation ev_agg
                                          WHERE deleted_at IS NOT NULL
                                          GROUP BY record_src) AS deleted_ev ON deleted_ev.record_src = event_agg.record_src
                      GROUP BY relay_prep.record_src, relay_prep.record_src
                      UNION ALL
                      SELECT 'neo4j',
                             _id   AS relay_id,
                             name AS relay_name,
                             NULL AS event_count,
                             NULL AS event_deleted_count,
                             NULL AS pubkey_count
                      FROM rdv.dump_neo4j_sat_cur
                      WHERE _labels = ':Relay'$$;
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relay_relay_id ON udm.Relay (relay_id)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_udm_relay_relay_name ON udm.Relay (relay_name)';
    ELSE
        BEGIN
            EXECUTE 'TRUNCATE TABLE udm.Relay';
            EXECUTE $$INSERT INTO udm.Relay (record_src, relay_id, relay_name, event_count, event_deleted_count, pubkey_count)
                          WITH relay_prep AS (SELECT DISTINCT record_src
                                              FROM udm.events_aggregation
                                              UNION
                                              SELECT DISTINCT record_src
                                              FROM udm.logs_aggregation)
                          SELECT relay_prep.record_src,
                                 NULL                             AS relay_id,
                                 relay_prep.record_src            AS relay_name,
                                 COUNT(event_agg.event_id)        AS event_count,
                                 COALESCE(MAX(deleted_ev.cnt), 0) AS event_deleted_count,
                                 COUNT(DISTINCT event_agg.pubkey) AS pubkey_count
                          FROM relay_prep
                                   LEFT JOIN udm.events_aggregation event_agg ON event_agg.record_src = relay_prep.record_src
                                   LEFT JOIN (SELECT record_src, COUNT(event_id) AS cnt
                                              FROM udm.events_aggregation ev_agg
                                              WHERE deleted_at IS NOT NULL
                                              GROUP BY record_src) AS deleted_ev ON deleted_ev.record_src = event_agg.record_src
                          GROUP BY relay_prep.record_src, relay_prep.record_src
                          UNION ALL
                          SELECT 'neo4j',
                                 _id   AS relay_id,
                                 name AS relay_name,
                                 NULL AS event_count,
                                 NULL AS event_deleted_count,
                                 NULL AS pubkey_count
                          FROM rdv.dump_neo4j_sat_cur
                          WHERE _labels = ':Relay'$$;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Error while refreshing udm.Relay: %', SQLERRM;
                RAISE;
        END;
    END IF;
END;
$proc$;
