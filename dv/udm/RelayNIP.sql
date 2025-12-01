/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Relationship table for relay to relay nip relationship
Modification:
=============================================
*/

DROP TABLE IF EXISTS udm.RelayNip;

SELECT *
INTO udm.RelayNip
FROM (SELECT 'neo4j' AS record_src,
             _start  AS relay_id,
             _end    AS nip_id
      FROM rdv.dump_neo4j_sat_cur
      WHERE _type = 'IMPLEMENTS') AS src;