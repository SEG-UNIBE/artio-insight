/*
=============================================
Author:				Michael Kaiser
Create date:		2026-04-01
Description:		View for aggregation of all the event tables from rdv including history
Modification:
=============================================
*/

DROP VIEW IF EXISTS udm.events_aggregation_history;

CREATE VIEW udm.events_aggregation_history AS
(
SELECT CAST('relay.artiostr.ch' AS VARCHAR(64)) AS record_src,
       id,
       created_at,
       updated_at,
       deleted_at,
       created,
       event_id,
       pubkey,
       kind,
       content,
       sig,
       tags
FROM rdv.events_artiostr_sat
UNION ALL
select CAST('relay.artio.inf.unibe.ch' AS VARCHAR(64)) AS record_src,
       id,
       created_at,
       updated_at,
       deleted_at,
       created,
       event_id,
       pubkey,
       kind,
       content,
       sig,
       tags
FROM rdv.events_artioinf_sat
    );
