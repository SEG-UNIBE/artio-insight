/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Table for getting stats about the user
Modification:
2025-11-25  mkaiser     Changing the script to physical table and including neo4j data
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.user_load(p_force boolean DEFAULT FALSE)
    LANGUAGE plpgsql
AS
$proc$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'udm') THEN
        RAISE EXCEPTION $$Schema udm does not exist$$;
    END IF;

    IF p_force THEN
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'user') THEN
            EXECUTE $$DROP TABLE udm.User$$;
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'user') THEN
        EXECUTE $$CREATE TABLE udm.User AS SELECT record_src, NULL AS user_id, pubkey, MIN(created_at) AS first_seen, MAX(updated_at) AS last_seen, COUNT(pubkey) AS seen_count, COUNT(event_id) AS event_count FROM udm.events_aggregation GROUP BY pubkey, record_src UNION ALL SELECT 'neo4j' AS record_src, _id AS user_id, pubkey, NULL AS first_seen, NULL AS last_seen, NULL AS seen_count, NULL AS event_count FROM rdv.dump_neo4j_sat_cur WHERE _labels = ':User'$$;
        EXECUTE $$CREATE INDEX IF NOT EXISTS idx_udm_user_pubkey ON udm.User (pubkey)$$;
        EXECUTE $$CREATE INDEX IF NOT EXISTS idx_udm_user_user_id ON udm.User (user_id)$$;
    ELSE
        BEGIN
            EXECUTE $$TRUNCATE TABLE udm.user$$;
            EXECUTE $$INSERT INTO udm.User (record_src, user_id, pubkey, first_seen, last_seen, seen_count, event_count) SELECT record_src, NULL AS user_id, pubkey, MIN(created_at) AS first_seen, MAX(updated_at) AS last_seen, COUNT(pubkey) AS seen_count, COUNT(event_id) AS event_count FROM udm.events_aggregation GROUP BY pubkey, record_src UNION ALL SELECT 'neo4j' AS record_src, _id AS user_id, pubkey, NULL AS first_seen, NULL AS last_seen, NULL AS seen_count, NULL AS event_count FROM rdv.dump_neo4j_sat_cur WHERE _labels = ':User'$$;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE $$Error while refreshing udm.User: %$$, SQLERRM;
                RAISE;
        END;
    END IF;

END;
$proc$;
