CREATE TABLE entries
(
    id           TEXT PRIMARY KEY,
    user_id      TEXT      NOT NULL,
    title        TEXT      NOT NULL,
    content      TEXT      NOT NULL,
    growth_stage TEXT      NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);