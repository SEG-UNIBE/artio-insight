/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Relationship table for relay to software usage
Modification:
=============================================
*/

DROP TABLE IF EXISTS udm.RelayOwner;

SELECT *
INTO udm.RelayOwner
FROM (SELECT 'neo4j' AS record_src,
             _start  AS user_id,
             _end    AS relay_id
      FROM rdv.dump_neo4j_sat_cur
      WHERE _type = 'OWNS') AS src;