/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Table for getting stats about the user
Modification:
2025-11-25  mkaiser     Changing the script to physical table and including neo4j data
=============================================
*/
DROP TABLE IF EXISTS udm.User;

SELECT *
INTO udm.user
FROM (SELECT record_src,
             NULL            AS user_id,
             pubkey,
             MIN(created_at) AS first_seen,
             MAX(updated_at) AS last_seen,
             COUNT(pubkey)   AS seen_count,
             COUNT(event_id) AS event_count

      FROM udm.events_aggregation
      GROUP BY pubkey, record_src
      UNION ALL
      SELECT 'neo4j' AS record_src,
             _id     AS user_id,
             pubkey,
             NULL    AS first_seen,
             NULL    AS last_seen,
             NULL    AS seen_count,
             NULL    AS event_count
      FROM rdv.dump_neo4j_sat_cur
      WHERE _labels = ':User') AS src;

