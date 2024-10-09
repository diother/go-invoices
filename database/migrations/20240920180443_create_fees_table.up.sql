CREATE TABLE fees (
    id TEXT PRIMARY KEY,
    description TEXT NOT NULL,
    created INTEGER NOT NULL,
    fee INTEGER NOT NULL,
    payout_id TEXT,
    FOREIGN KEY (payout_id) REFERENCES payouts(id)
);

CREATE INDEX idx_fees_payout_id ON fees (payout_id);
