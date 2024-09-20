CREATE TABLE payouts (
    id VARCHAR(255) PRIMARY KEY,
    created BIGINT UNSIGNED NOT NULL,
    gross INT UNSIGNED NOT NULL,
    fee INT UNSIGNED NOT NULL,
    net INT UNSIGNED NOT NULL,
    INDEX idx_created (created)
);
