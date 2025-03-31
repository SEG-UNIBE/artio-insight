CREATE TABLE IF NOT EXISTS current_nip_tags (
    nip VARCHAR(255) PRIMARY KEY NOT NULL,
    final BOOLEAN NOT NULL DEFAULT false,
    draft BOOLEAN NOT NULL DEFAULT false,
    mandatory BOOLEAN NOT NULL DEFAULT false,
    optional BOOLEAN NOT NULL DEFAULT false,
    recommended BOOLEAN NOT NULL DEFAULT false,
    unrecommended BOOLEAN NOT NULL DEFAULT false,
    last_updated DATE NOT NULL
);

