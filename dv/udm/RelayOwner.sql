/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Relationship table for relay to owner usage
Modification:
=============================================
*/

DROP TABLE IF EXISTS udm.RelayOwner;

SELECT *
INTO udm.RelaySoftware
FROM (SELECT 'neo4j' AS record_src,
             _start  AS relay_id,
             _end    AS software_id
      FROM rdv.dump_neo4j_sat_cur
      WHERE _type = 'USES_SOFTWARE') AS src;