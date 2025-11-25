/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Table representing the software dimension
Modification:
=============================================
*/

DROP TABLE IF EXISTS udm.Software;

SELECT *
INTO udm.Software
FROM (SELECT 'neo4j'  AS record_src,
             _id      AS software_id,
             software AS software
      FROM rdv.dump_neo4j_sat_cur
      WHERE _labels = ':Software') AS src;
