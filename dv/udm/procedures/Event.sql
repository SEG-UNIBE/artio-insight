/*
=============================================
Author:				Michael Kaiser
Create date:		2025-10-19
Description:		Table for all events
Modification:
2025-11-25  mkaiser     Changing the script to physical table
2026-01-12  mkaiser     Refactoring to procedure instead of loading script
=============================================
*/

CREATE OR REPLACE PROCEDURE udm.event_load(p_force boolean DEFAULT FALSE)
    LANGUAGE plpgsql
AS
$proc$
BEGIN
    IF p_force THEN
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'event') THEN
            EXECUTE $$DROP TABLE udm.Event$$;
        END IF;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'udm' AND table_name = 'event') THEN
        EXECUTE $$CREATE TABLE udm.Event AS
                    SELECT record_src, id,
                    created_at,
                    updated_at,
                    deleted_at,
                    created,
                    event_id,
                    pubkey,
                    kind,
                    content,
                    sig,
                    tags
                    FROM udm.events_aggregation$$;
        EXECUTE $$CREATE INDEX IF NOT EXISTS idx_udm_event_event_id ON udm.Event (event_id)$$;
        EXECUTE $$CREATE INDEX IF NOT EXISTS idx_udm_event_pubkey ON udm.Event (pubkey)$$;
    ELSE
        -- Use a transaction-safe refresh
        BEGIN
            EXECUTE $$TRUNCATE TABLE udm.Event$$;
            EXECUTE $$INSERT INTO udm.Event (record_src, id, created_at, updated_at, deleted_at, created, event_id, pubkey, kind, content, sig, tags) SELECT record_src, id, created_at, updated_at, deleted_at, created, event_id, pubkey, kind, content, sig, tags FROM udm.events_aggregation$$;
        EXCEPTION
            WHEN OTHERS THEN
                -- Re-raise with context
                RAISE NOTICE $$Error while refreshing udm.Event: %$$, SQLERRM;
                RAISE;
        END;
    END IF;

END;
$proc$;

