/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		View for all realtionships from the neo4j datasource
Modification:
=============================================
*/

DROP VIEW IF EXISTS udm.neo4j_relationships;

CREATE VIEW udm.neo4j_relationships AS
(
SELECT _start AS source_node_id,
       _end   AS target_node_id,
       _type  AS type
FROM rdv.dump_neo4j_sat_cur
WHERE _type <> ''
    );
