-- +goose Up
-- +goose StatementBegin

CREATE TABLE tenable_settings (
    id INT IDENTITY(1,1) PRIMARY KEY,
    base_url VARCHAR(255) NOT NULL,
    access_key VARCHAR(300) NOT NULL,
    secret_key VARCHAR(300) NOT NULL,
    sync_frequency_hours INT NOT NULL DEFAULT 24,      
    lookback_days INT NOT NULL DEFAULT 30,             
    status VARCHAR(10) NOT NULL 
        CONSTRAINT chk_tenable_status CHECK (status IN ('active', 'inactive'))
        DEFAULT 'active',
    created_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME(),
    updated_at DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME()
);


-- default row
INSERT INTO tenable_settings 
    (base_url, access_key, secret_key, sync_frequency_hours, lookback_days, status)
VALUES
    ('https://sc.tenalab.online', '', '', 24, 30, 'active');

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE tenable_settings;

-- +goose StatementEnd
