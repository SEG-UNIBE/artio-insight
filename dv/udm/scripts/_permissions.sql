/*
=============================================
Author:				Michael Kaiser
Create date:		2026-01-09
Description:		script to reset the permissions after udm load
Modification:
=============================================
*/

GRANT SELECT ON ALL TABLES IN SCHEMA udm TO udm_reader;
GRANT USAGE ON SCHEMA udm TO udm_reader;
ALTER DEFAULT PRIVILEGES IN SCHEMA udm GRANT SELECT ON TABLES TO udm_reader;


GRANT ALL ON ALL TABLES IN SCHEMA udm TO udm_admin;
GRANT USAGE ON SCHEMA udm TO udm_admin;
ALTER DEFAULT PRIVILEGES IN SCHEMA udm GRANT ALL ON TABLES TO udm_admin;
