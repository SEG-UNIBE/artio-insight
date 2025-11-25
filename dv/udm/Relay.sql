/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Table for getting stats about the relays
Modification:
2025-11-25  mkaiser     Changing the script to physical table and including neo4j data
=============================================
*/
DROP TABLE IF EXISTS udm.Relay;


WITH relay_prep AS (SELECT DISTINCT record_src
                    FROM udm.events_aggregation
                    UNION
                    SELECT DISTINCT record_src
                    FROM udm.logs_aggregation)
SELECT *
INTO udm.Relay
FROM (SELECT relay_prep.record_src,
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
      WHERE _labels = ':Relay') AS src;