/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Relationship table for relay to relay detection usage
Modification:
=============================================
*/

DROP TABLE IF EXISTS udm.RelayDetection;

SELECT *
INTO udm.RelayDetection
FROM (SELECT 'neo4j' AS record_src,
             _start  AS relay_id,
             _end    AS detected_relay_id
      FROM rdv.dump_neo4j_sat_cur
      WHERE _type = 'DETECTED') AS src;