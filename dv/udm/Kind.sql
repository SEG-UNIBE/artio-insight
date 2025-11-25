/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Table for the kinds of events
Modification:
2025-11-25  mkaiser     Changing the script to physical table
=============================================
*/
DROP TABLE IF EXISTS udm.kind;

SELECT *
INTO udm.kind
FROM (SELECT record_src,
             kind,
             COUNT(DISTINCT pubkey) AS user_count,
             COUNT(event_id)        AS event_count

      FROM udm.events_aggregation
      GROUP BY kind, record_src) AS src;
