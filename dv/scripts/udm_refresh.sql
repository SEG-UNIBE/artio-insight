/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Script to refresh all materialized views in the udm schema
Modification:
=============================================
*/

DO
$$
    DECLARE
        output        text := '';
    BEGIN
        -- loading tables to load
        DROP TABLE IF EXISTS temp_tables;
        CREATE TEMP TABLE temp_tables AS (select * from pg_matviews where schemaname = 'udm');

        DECLARE
            curs CURSOR FOR SELECT "matviewname"
                            FROM temp_tables;
            tab_name TEXT;
        BEGIN
            OPEN curs;
            LOOP
                FETCH curs INTO tab_name;
                EXIT WHEN NOT FOUND;

                SELECT INTO output CONCAT(output, FORMAT(E'REFRESH MATERIALIZED VIEW udm.%s;\n', tab_name));

            END LOOP;
            CLOSE curs;
        END;

        -- show or execute the script
        RAISE NOTICE '%', output;
        RAISE NOTICE 'executing udm materialized view refresh script';
        EXECUTE output;
    END;
$$;
