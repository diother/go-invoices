CREATE TABLE payouts (
    id TEXT PRIMARY KEY,
    created INTEGER NOT NULL,
    gross INTEGER NOT NULL,
    fee INTEGER NOT NULL,
    net INTEGER NOT NULL
);

CREATE INDEX idx_created ON payouts (created);
