-- +goose Up
-- +goose StatementBegin
CREATE TABLE system_preferences (
    id INT PRIMARY KEY IDENTITY(1,1),
    currency_id INT NOT NULL,
    timezone_id INT NOT NULL,
    date_format VARCHAR(45) NOT NULL,
    time_format VARCHAR(10) NOT NULL DEFAULT '12',
        CONSTRAINT CHK_system_preferences_time_format CHECK (time_format IN ('12', '24')),
    doc_upload_size INT NOT NULL,
    logo_upload_size INT NOT NULL,
    idle_timeout INT NOT NULL,
    is_maintenance_mode_on TINYINT DEFAULT 0,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_system_preferences_currencies FOREIGN KEY (currency_id) REFERENCES currencies (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_preferences_timezones FOREIGN KEY (timezone_id) REFERENCES timezones (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_preferences_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_preferences_users_3 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS system_preferences;
-- +goose StatementEnd
