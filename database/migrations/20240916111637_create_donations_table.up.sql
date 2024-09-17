CREATE TABLE donations (
    id VARCHAR(255) PRIMARY KEY,
    created BIGINT UNSIGNED,
    gross INT UNSIGNED,
    fee INT UNSIGNED,
    net INT UNSIGNED,
    product VARCHAR(255),
    client_name VARCHAR(255),
    client_email VARCHAR(255),
    payout_id VARCHAR(255),
    FOREIGN KEY (payout_id) REFERENCES payouts(id),
    INDEX idx_created (created),
    INDEX idx_payout_id (payout_id)
);
