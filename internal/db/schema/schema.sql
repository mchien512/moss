CREATE TABLE entries (
                         id TEXT PRIMARY KEY,
                         user_id TEXT NOT NULL,
                         title TEXT NOT NULL,
                         content TEXT NOT NULL,
                         mood TEXT NOT NULL,
                         created_at TIMESTAMPTZ NOT NULL,
                         updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX entries_user_id_idx ON entries(user_id);
CREATE INDEX entries_updated_at_idx ON entries(updated_at);
