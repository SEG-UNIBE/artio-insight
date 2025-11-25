/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Table representing the RelayAlternativeName dimension
Modification:
=============================================
*/

DROP TABLE IF EXISTS udm.RelayAlternativeName;

SELECT *
INTO udm.RelayAlternativeName
FROM (SELECT 'neo4j' AS record_src,
             _id     AS relayaltname_id,
             name    AS relayaltname
      FROM rdv.dump_neo4j_sat_cur
      WHERE _labels = ':RelayAlternativeName') AS src;


