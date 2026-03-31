/*
=============================================
Author:				Michael Kaiser
Create date:		2026-03-18
Description:		View for relay with additional information for easy and fast retreival
Modification:
=============================================
*/

DROP VIEW IF EXISTS udm.v_relay;

CREATE VIEW udm.v_relay AS
(
SELECT r.relay_id,
       r.relay_name,
       r.is_valid,
       r.valid_reason,
       s.software_id,
       s.software,
       o.user_id       AS owner_user_id,
       o.pubkey        AS owner_pubkey,
       r2.relay_id     AS detected_relay_id,
       r2.relay_name   AS detected_relay_name,
       r2.is_valid     AS detected_relay_is_valid,
       r2.valid_reason AS detected_relay_valid_reason,
       s2.software_id  AS detected_software_id,
       s2.software     AS detected_software
FROM udm.relay r
         LEFT JOIN udm.relaysoftware rs
                   ON r.relay_id = rs.relay_id
         LEFT JOIN udm.software s
                   ON s.software_id = rs.software_id
         LEFT JOIN udm.relayowner ro
                   ON ro.relay_id = r.relay_id
         LEFT JOIN udm."user" o
                   ON o.user_id = ro.user_id
         LEFT JOIN udm.relaydetection rd
                   ON rd.relay_id = r.relay_id
         LEFT JOIN udm.relay r2
                   ON rd.detected_relay_id = r2.relay_id
         LEFT JOIN udm.relaysoftware rs2
                   ON r2.relay_id = rs2.relay_id
         LEFT JOIN udm.software s2
                   ON s2.software_id = rs2.software_id
    );
