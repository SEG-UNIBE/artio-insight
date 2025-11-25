/*
=============================================
Author:				Michael Kaiser
Create date:		2025-11-25
Description:		Relationship table for relay to relayaltname usage
Modification:
=============================================
*/

DROP TABLE IF EXISTS udm.RelayRelayAltName;

SELECT *
INTO udm.RelayRelayAltName
FROM (SELECT 'neo4j' AS record_src,
             _start  AS relay_id,
             _end    AS relayaltname_id
      FROM rdv.dump_neo4j_sat_cur
      WHERE _type = 'ALT_NAME') AS src;