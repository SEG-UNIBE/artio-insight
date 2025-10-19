/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Materialized view for all events
Modification:
=============================================
*/
DROP MATERIALIZED VIEW IF EXISTS udm.Event;

CREATE MATERIALIZED VIEW udm.Event AS
(
SELECT record_src,
       id,
       created_at,
       updated_at,
       deleted_at,
       created,
       event_id,
       pubkey,
       kind,
       content,
       sig,
       tags
FROM udm.events_aggregation
);
