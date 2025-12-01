/*
=============================================
Author:				Michael Kaiser
Create date:		2025-12-01
Description:		Table representing the NIP-Dimension
Modification:
=============================================
*/

DROP TABLE IF EXISTS udm.Nip;

SELECT *
INTO udm.Nip
FROM (SELECT 'neo4j'              AS record_src,
             _id                  AS nip_id,
             name                 AS nip_name,
             CONCAT('NIP-', name) AS nip_full_name
      FROM rdv.dump_neo4j_sat_cur
      WHERE _labels = ':NIP') AS src;
