/*
=============================================
Author:				Michael Kaiser
Create date:		2025-08-20
Description:		Script to setup the needed database schemas and other structures for the insight DV.
Modification:
YYYY-MM-DD	mkaiser
=============================================
*/

DROP SCHEMA IF EXISTS cfg;
DROP SCHEMA IF EXISTS inb;
DROP SCHEMA IF EXISTS rdv;
DROP SCHEMA IF EXISTS udm;

CREATE SCHEMA cfg;
CREATE SCHEMA inb;
CREATE SCHEMA rdv;
CREATE SCHEMA udm;


/*
BEGIN:              cfg
 */
DROP OWNED BY cfg_reader;
DROP ROLE IF EXISTS cfg_reader;
CREATE ROLE cfg_reader WITH
    NOLOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    INHERIT
    NOREPLICATION
    NOBYPASSRLS
    CONNECTION LIMIT -1;
COMMENT ON ROLE cfg_reader IS 'Role with read-only permission on the cfg schema';

GRANT SELECT ON ALL TABLES IN SCHEMA cfg TO cfg_reader;
GRANT USAGE ON SCHEMA cfg TO cfg_reader;

DROP OWNED BY cfg_admin;
DROP ROLE IF EXISTS cfg_admin;
CREATE ROLE cfg_admin WITH
    NOLOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    INHERIT
    NOREPLICATION
    NOBYPASSRLS
    CONNECTION LIMIT -1;
COMMENT ON ROLE cfg_admin IS 'Role with administrative permission on the cfg schema';

GRANT ALL ON ALL TABLES IN SCHEMA cfg TO cfg_admin;
GRANT USAGE ON SCHEMA cfg TO cfg_admin;
/*
END:              cfg
 */

/*
BEGIN:              inb
*/
DROP OWNED BY inb_reader;
DROP ROLE IF EXISTS inb_reader;
CREATE ROLE inb_reader WITH
    NOLOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    INHERIT
    NOREPLICATION
    NOBYPASSRLS
    CONNECTION LIMIT -1;
COMMENT ON ROLE inb_reader IS 'Role with read-only permission on the inb schema';

GRANT SELECT ON ALL TABLES IN SCHEMA inb TO inb_reader;
GRANT USAGE ON SCHEMA inb TO inb_reader;

DROP OWNED BY inb_admin;
DROP ROLE IF EXISTS inb_admin;
CREATE ROLE inb_admin WITH
    NOLOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    INHERIT
    NOREPLICATION
    NOBYPASSRLS
    CONNECTION LIMIT -1;
COMMENT ON ROLE inb_admin IS 'Role with administrative permission on the inb schema';

GRANT ALL ON ALL TABLES IN SCHEMA inb TO inb_admin;
GRANT USAGE ON SCHEMA inb TO inb_admin;

/*
END:              inb
 */

/*
BEGIN:              rdv
 */
DROP OWNED BY rdv_reader;
DROP ROLE IF EXISTS rdv_reader;
CREATE ROLE rdv_reader WITH
    NOLOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    INHERIT
    NOREPLICATION
    NOBYPASSRLS
    CONNECTION LIMIT -1;
COMMENT ON ROLE rdv_reader IS 'Role with read-only permission on the rdv schema';

GRANT SELECT ON ALL TABLES IN SCHEMA rdv TO rdv_reader;
GRANT USAGE ON SCHEMA rdv TO rdv_reader;

DROP OWNED BY rdv_admin;
DROP ROLE IF EXISTS rdv_admin;
CREATE ROLE rdv_admin WITH
    NOLOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    INHERIT
    NOREPLICATION
    NOBYPASSRLS
    CONNECTION LIMIT -1;
COMMENT ON ROLE rdv_admin IS 'Role with administrative permission on the rdv schema';

GRANT ALL ON ALL TABLES IN SCHEMA rdv TO rdv_admin;
GRANT USAGE ON SCHEMA rdv TO rdv_admin;
/*
END:              rdv
 */

/*
BEGIN:              udm
*/
DROP OWNED BY udm_reader;
DROP ROLE IF EXISTS udm_reader;
CREATE ROLE udm_reader WITH
    NOLOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    INHERIT
    NOREPLICATION
    NOBYPASSRLS
    CONNECTION LIMIT -1;
COMMENT ON ROLE udm_reader IS 'Role with read-only permission on the udm schema';

GRANT SELECT ON ALL TABLES IN SCHEMA udm TO udm_reader;
GRANT USAGE ON SCHEMA udm TO udm_reader;

DROP OWNED BY udm_admin;
DROP ROLE IF EXISTS udm_admin;
CREATE ROLE udm_admin WITH
    NOLOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    INHERIT
    NOREPLICATION
    NOBYPASSRLS
    CONNECTION LIMIT -1;
COMMENT ON ROLE udm_admin IS 'Role with administrative permission on the udm schema';

GRANT ALL ON ALL TABLES IN SCHEMA udm TO udm_admin;
GRANT USAGE ON SCHEMA udm TO udm_admin;

/*
END:              udm
 */

ALTER DATABASE insight SET search_path TO inb;