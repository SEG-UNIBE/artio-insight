/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Table for all events
Modification:
2025-11-25  mkaiser     Changing the script to physical table
=============================================
*/
DROP TABLE IF EXISTS udm.Event;

SELECT *
INTO udm.Event
FROM (SELECT record_src,
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
      FROM udm.events_aggregation) AS src;
