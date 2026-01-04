/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Relationship table for relay to software usage
Modification:
2026-01-04  mkaiser     Fix the mix of software and owner
=============================================
*/

DROP TABLE IF EXISTS udm.RelaySoftware;

SELECT *
INTO udm.RelaySoftware
FROM (SELECT 'neo4j' AS record_src,
             _start  AS relay_id,
             _end    AS software_id
      FROM rdv.dump_neo4j_sat_cur
      WHERE _type = 'USES_SOFTWARE') AS src;