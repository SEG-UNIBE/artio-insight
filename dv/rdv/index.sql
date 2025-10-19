/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Additional index creation in rdv
Modification:
=============================================
*/

CREATE INDEX IF NOT EXISTS idx_logs_artiostr_sat_cur ON rdv.logs_artiostr_sat_cur (id, ip, type);
CREATE INDEX IF NOT EXISTS idx_events_artiostr_sat_cur ON rdv.events_artiostr_sat_cur (event_id, pubkey, kind);

CREATE INDEX IF NOT EXISTS idx_logs_artioinf_sat_cur ON rdv.logs_artioinf_sat_cur (id, ip, type);
CREATE INDEX IF NOT EXISTS idx_events_artioinf_sat_cur ON rdv.events_artioinf_sat_cur (event_id, pubkey, kind);