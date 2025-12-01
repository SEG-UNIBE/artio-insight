/*
=============================================
Author:				Michael Kaiser
Create date:		2025-12-01
Description:		Relationship table for relay to nip kinds relationship
Modification:
=============================================
*/

DROP TABLE IF EXISTS udm.NipKind;

SELECT *
INTO udm.NipKind
FROM (SELECT 'nostrbook.dev' AS record_src,
             nip_name  AS nip_name,
             kind    AS kind
      FROM rdv.nipkinds_nipfile_sat_cur) AS src;