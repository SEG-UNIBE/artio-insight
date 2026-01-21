/*
=============================================
Author:				Michael Kaiser
Create date:		2026-01-21
Description:		Procedure creation for running all udm procedures in the database
Modification:
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.run_all_loads(p_force boolean DEFAULT FALSE)
    LANGUAGE plpgsql
AS
$proc$
DECLARE
    r            RECORD;
    call_sql     TEXT;
    msg          TEXT;
    called_count INT := 0;
    failed_count INT := 0;
BEGIN
    -- Verify schema exists
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'udm') THEN
        RAISE EXCEPTION 'Schema udm does not exist';
    END IF;

    FOR r IN
        SELECT routine_name
        FROM information_schema.routines
        WHERE routine_schema = 'udm'
          AND routine_type = 'PROCEDURE'
          AND routine_name <> 'run_all_loads'
        ORDER BY routine_name
        LOOP
            -- Build CALL SQL, pass boolean literal
            call_sql := FORMAT('CALL udm.%I(%s)', r.routine_name, CASE WHEN p_force THEN 'true' ELSE 'false' END);
            RAISE INFO 'Calling %', call_sql;
            BEGIN
                EXECUTE call_sql;
                called_count := called_count + 1;
                RAISE INFO $$Succeeded: %$$, r.routine_name;
            EXCEPTION
                WHEN OTHERS THEN
                    failed_count := failed_count + 1;
                    msg := FORMAT('Procedure %s failed: %s', r.routine_name, SQLERRM);
                    RAISE INFO $$%$$, msg;
            END;
        END LOOP;

    IF p_force THEN
        GRANT SELECT ON ALL TABLES IN SCHEMA udm TO udm_reader;
        GRANT USAGE ON SCHEMA udm TO udm_reader;

        GRANT ALL ON ALL TABLES IN SCHEMA udm TO udm_admin;
        GRANT USAGE ON SCHEMA udm TO udm_admin;
    END IF;

    msg := FORMAT('udm.run_all_loads finished: called=%s called, failed=%s failed', called_count, failed_count);
    RAISE INFO $$%$$, msg;
END;
$proc$;
