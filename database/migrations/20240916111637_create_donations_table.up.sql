CREATE TABLE donations (
    id TEXT PRIMARY KEY,
    created INTEGER NOT NULL,
    gross INTEGER NOT NULL,
    fee INTEGER NOT NULL,
    net INTEGER NOT NULL,
    client_name TEXT NOT NULL,
    client_email TEXT NOT NULL,
    payout_id TEXT,
    FOREIGN KEY (payout_id) REFERENCES payouts(id)
);

CREATE INDEX idx_donations_payout_id ON donations (payout_id);
