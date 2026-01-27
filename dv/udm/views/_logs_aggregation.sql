/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		View for aggregation of all the event tables from rdv
Modification:
=============================================
*/

DROP VIEW IF EXISTS udm.logs_aggregation;

CREATE VIEW udm.logs_aggregation AS
(
SELECT CAST('relay.artiostr.ch' AS VARCHAR(64)) AS record_src,
       id,
       created_at,
       updated_at,
       deleted_at,
       ip,
       type,
       content
FROM rdv.logs_artiostr_sat_cur
UNION ALL
select CAST('relay.artio.inf.unibe.ch' AS VARCHAR(64)) AS record_src,
       id,
       created_at,
       updated_at,
       deleted_at,
       ip,
       type,
       content
FROM rdv.logs_artioinf_sat_cur
    );
