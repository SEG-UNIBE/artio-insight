/*
=============================================
Author:				Michael Kaiser
Create date:		2025-08-20
Description:		Script to generate and possibly execute the load for the rdv sat tables
Modification:
YYYY-MM-DD	mkaiser
=============================================
*/

DO
$$
    DECLARE
        executeScript bool := TRUE;
        output        text := '';
        tableCount    int  := 0;
    BEGIN
        -- loading tables to load
        DROP TABLE IF EXISTS temp_tables;
        CREATE TEMP TABLE temp_tables AS (SELECT "table_name"
                                          FROM INFORMATION_SCHEMA.TABLES
                                          WHERE TABLE_SCHEMA = 'inb');

        SELECT INTO tableCount COUNT(*) FROM temp_tables;
        RAISE NOTICE 'Tables found in inb layer: %', tableCount;

        DECLARE
            curs CURSOR FOR SELECT "table_name"
                            FROM temp_tables;
            tab_name TEXT;
        BEGIN
            OPEN curs;
            LOOP
                FETCH curs INTO tab_name;
                EXIT WHEN NOT FOUND;

                --SELECT INTO output CONCAT(output, E'BEGIN WORK;\n');
                --SELECT INTO output CONCAT(output, FORMAT(E'LOCK TABLE rdv.%s_sat_cur IN EXCLUSIVE MODE;\n', tab_name));
                SELECT INTO output CONCAT(output, FORMAT(E'UPDATE rdv.%s_sat ' ||
                                                         E'SET delete_dts = NOW() ' ||
                                                         E'WHERE frh NOT IN (SELECT frh FROM rdv.%s_sat_cur) ' ||
                                                         E'AND delete_dts IS NOT NULL;\n', tab_name, tab_name));
                SELECT INTO output CONCAT(output, FORMAT(E'INSERT INTO rdv.%s_sat ' ||
                                                         E'SELECT * FROM rdv.%s_sat_cur ' ||
                                                         E'WHERE frh NOT IN (SELECT frh FROM rdv.%s_sat);\n \n', tab_name, tab_name, tab_name));
                --SELECT INTO output CONCAT(output, E'COMMIT WORK;\n \n');

            END LOOP;
            CLOSE curs;
        END;

        -- show or execute the script
        RAISE NOTICE '%', output;
        IF executeScript THEN
            RAISE NOTICE 'executing rdv table load script';
            EXECUTE output;
        END IF;
    END;
$$;
