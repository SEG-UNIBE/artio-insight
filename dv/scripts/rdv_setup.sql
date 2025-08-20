/*
=============================================
Author:				Michael Kaiser
Create date:		2025-08-20
Description:		Script to setup the tables in the rdv based on the structure of the tables in the rdv layer layer
Modification:
YYYY-MM-DD	mkaiser
=============================================
*/

DO
$$
    DECLARE
        redeploy      bool := TRUE;
        executeScript bool := FALSE;
        output        text := '';
        tableCount    int  := 0;
    BEGIN
        -- loading tables to build
        DROP TABLE IF EXISTS temp_tables;
        CREATE TEMP TABLE temp_tables AS (SELECT "table_name"
                                          FROM INFORMATION_SCHEMA.TABLES
                                          WHERE TABLE_SCHEMA = 'inb');

        SELECT INTO tableCount COUNT(*) FROM temp_tables;
        RAISE NOTICE 'Tables found in inb layer: %', tableCount;

        DECLARE
            curs CURSOR FOR SELECT "table_name"
                            FROM temp_tables;
            tab_name   TEXT;
            tableCount integer := 0;
        BEGIN
            OPEN curs;
            LOOP
                FETCH curs INTO tab_name;
                EXIT WHEN NOT FOUND;

                SELECT INTO tableCount COUNT(*)
                FROM INFORMATION_SCHEMA.TABLES
                WHERE TABLE_SCHEMA = 'rdv'
                  AND table_name IN (CONCAT(tab_name, '_sat'), CONCAT(tab_name, '_sat_cur'));

                IF tableCount = 2 AND redeploy IS FALSE THEN
                    RAISE NOTICE 'RDV Tables already exist for table: %s', tab_name;
                    CONTINUE;
                END IF;

                SELECT INTO output CONCAT(output, FORMAT(E'DROP TABLE IF EXISTS rdv.%s_sat;\n', tab_name));
                SELECT INTO output CONCAT(output, FORMAT(
                        E'CREATE TABLE rdv.%s_sat AS (SELECT CAST(NULL AS timestamp with time zone) AS load_dts, CAST(NULL AS timestamp with time zone) AS delete_dts, * FROM inb.%s WHERE 1=2);\n',
                        tab_name, tab_name));

                SELECT INTO output CONCAT(output, FORMAT(E'DROP TABLE IF EXISTS rdv.%s_sat_cur;\n', tab_name));
                SELECT INTO output CONCAT(output, FORMAT(
                        E'CREATE TABLE rdv.%s_sat_cur AS (SELECT CAST(NULL AS timestamp with time zone) AS load_dts, CAST(NULL AS timestamp with time zone) AS delete_dts, * FROM inb.%s WHERE 1=2);\n',
                        tab_name, tab_name));

            END LOOP;
            CLOSE curs;
        END;

        -- show or execute the script
        RAISE NOTICE '%', output;
        IF executeScript THEN
            RAISE NOTICE 'executing rdv table creation script';
            EXECUTE output;
        END IF;
    END;
$$;