CREATE TABLE entry_links (
     source_entry_id TEXT,
     target_entry_id TEXT,
     user_id TEXT NOT NULL,
     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

     PRIMARY KEY (source_entry_id, target_entry_id),
     FOREIGN KEY (source_entry_id) REFERENCES entries(id),
     FOREIGN KEY (target_entry_id) REFERENCES entries(id)
);

CREATE INDEX source_idx ON entry_links(source_entry_id);

CREATE INDEX target_idx ON entry_links(target_entry_id);

CREATE INDEX user_idx ON entry_links(user_id);
