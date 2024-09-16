CREATE TABLE payouts (
    id VARCHAR(255) PRIMARY KEY,
    created BIGINT UNSIGNED,
    gross INT UNSIGNED,
    fee INT UNSIGNED,
    net INT UNSIGNED,
    INDEX idx_created (created)
);
