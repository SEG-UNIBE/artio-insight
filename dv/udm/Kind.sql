/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Materialized view for the kinds of events
Modification:
=============================================
*/
DROP MATERIALIZED VIEW IF EXISTS udm.kind;

CREATE MATERIALIZED VIEW udm.kind AS
(
SELECT record_src,
       kind,
       COUNT(DISTINCT pubkey) AS user_count,
       COUNT(event_id)        AS event_count

FROM udm.events_aggregation
GROUP BY kind, record_src
    );
