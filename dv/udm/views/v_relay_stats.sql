/*
=============================================
Author:				Michael Kaiser
Create date:		2026-04-15
Description:		View for relay with basic stats and now discovery information for easy and fast retreival
Modification:
=============================================
*/

DROP VIEW IF EXISTS udm.v_relay_stats;

CREATE VIEW udm.v_relay_stats AS
(
SELECT r.relay_id,
       r.relay_name,
       r.is_valid,
       ip.ip_id,
       ip.ip_address,
       ip.ip_type,
       nip.nip_name,
       nip.nip_full_name,
       s.software
FROM udm.relay AS r
         LEFT JOIN udm.relayip AS rip
                   ON r.relay_id = rip.relay_id
         LEFT JOIN udm.ip AS ip
                   ON rip.ip_id = ip.ip_id
         LEFT JOIN udm.relaynip AS rnip
                   ON r.relay_id = rnip.relay_id
         LEFT JOIN udm.nip AS nip
                   ON rnip.nip_id = nip.nip_id
         LEFT JOIN udm.relaysoftware rs
                   ON r.relay_id = rs.relay_id
         LEFT JOIN udm.software s
                   ON s.software_id = rs.software_id
    );
