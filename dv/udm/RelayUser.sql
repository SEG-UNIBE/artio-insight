/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Relationship table for relay to user usage
Modification:
=============================================
*/

DROP TABLE IF EXISTS udm.RelayUser;

SELECT *
INTO udm.RelayUser
FROM (SELECT 'neo4j' AS record_src,
             _start  AS user_id,
             _end    AS relay_id
      FROM rdv.dump_neo4j_sat_cur
      WHERE _type = 'USES') AS src;