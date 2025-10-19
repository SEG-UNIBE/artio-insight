/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Materialized view for getting stats about the relays
Modification:
=============================================
*/
DROP MATERIALIZED VIEW IF EXISTS udm.Relay;

CREATE MATERIALIZED VIEW udm.Relay AS
(
WITH relay_temp AS (SELECT DISTINCT record_src
                    FROM udm.events_aggregation
                    UNION
                    SELECT DISTINCT record_src
                    FROM udm.logs_aggregation)

SELECT relay_temp.record_src,
       relay_temp.record_src            AS relay_name,
       COUNT(event_agg.event_id)        AS event_count,
       COALESCE(MAX(deleted_ev.cnt), 0) AS event_deleted_count,
       COUNT(DISTINCT event_agg.pubkey) AS pubkey_count
FROM relay_temp
         LEFT JOIN udm.events_aggregation event_agg ON event_agg.record_src = relay_temp.record_src
         LEFT JOIN (SELECT record_src, COUNT(event_id) AS cnt
                    FROM udm.events_aggregation ev_agg
                    WHERE deleted_at IS NOT NULL
                    GROUP BY record_src) AS deleted_ev ON deleted_ev.record_src = event_agg.record_src
GROUP BY relay_temp.record_src, relay_temp.record_src
);
