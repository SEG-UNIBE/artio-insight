/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Materialized view for getting stats about the user
Modification:
=============================================
*/
DROP MATERIALIZED VIEW IF EXISTS udm.User;

CREATE MATERIALIZED VIEW udm.User AS
(
SELECT record_src,
       pubkey,
       MIN(created_at) AS first_seen,
       MAX(updated_at) AS last_seen,
       COUNT(pubkey)   AS seen_count,
       COUNT(event_id) AS event_count

FROM udm.events_aggregation
GROUP BY pubkey, record_src
    );
