/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		View for all nodes from the neo4j datasource
Modification:
=============================================
*/

DROP VIEW IF EXISTS udm.neo4j_nodes;

CREATE VIEW udm.neo4j_nodes AS
(
SELECT _id AS node_id,
       _labels,
       isvalid,
       name,
       pubkey,
       software,
       validreason
FROM rdv.dump_neo4j_sat_cur
WHERE _labels <> ''
    );
