/*
=============================================
Author:				Michael Kaiser
Create date:		2026-01-09
Description:		Table representing the IP dimension
Modification:
=============================================
*/

DROP TABLE IF EXISTS udm.IP;

SELECT *
INTO udm.IP
FROM (SELECT 'neo4j'  AS record_src,
             _id      AS ip_id,
             address AS ip_address,
             CAST(address AS inet) as ip,
             family(CAST(address AS inet)) as ip_type
      FROM rdv.dump_neo4j_sat_cur
      WHERE _labels = ':IP') AS src;