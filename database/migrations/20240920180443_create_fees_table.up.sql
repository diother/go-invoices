CREATE TABLE fees (
    id VARCHAR(255) PRIMARY KEY,
    created BIGINT UNSIGNED NOT NULL,
    fee INT UNSIGNED NOT NULL,
    payout_id VARCHAR(255),
    FOREIGN KEY (payout_id) REFERENCES payouts(id),
    INDEX idx_payout_id (payout_id)
);
