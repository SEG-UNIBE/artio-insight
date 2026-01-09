/*
=============================================
Author:				Michael Kaiser
Create date:		2026-01-09
Description:		Relationship table for relay to relay ip relationship
Modification:
=============================================
*/

DROP TABLE IF EXISTS udm.RelayIp;

SELECT *
INTO udm.RelayIp
FROM (SELECT 'neo4j' AS record_src,
             _start  AS relay_id,
             _end    AS ip_id
      FROM rdv.dump_neo4j_sat_cur
      WHERE _type = 'HAS_IP') AS src;