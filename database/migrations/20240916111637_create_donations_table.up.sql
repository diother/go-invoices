CREATE TABLE donations (
    id VARCHAR(255) PRIMARY KEY,
    created BIGINT UNSIGNED NOT NULL,
    gross INT UNSIGNED NOT NULL,
    fee INT UNSIGNED NOT NULL,
    net INT UNSIGNED NOT NULL,
    client_name VARCHAR(255) NOT NULL,
    client_email VARCHAR(255) NOT NULL,
    payout_id VARCHAR(255),
    FOREIGN KEY (payout_id) REFERENCES payouts(id),
    INDEX idx_created (created),
    INDEX idx_payout_id (payout_id)
);
