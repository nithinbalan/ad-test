-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_settings (
    id INT PRIMARY KEY IDENTITY(1,1),
    user_id INT NOT NULL,
    timezone_id INT NOT NULL,
    language VARCHAR(2) NOT NULL DEFAULT 'en' 
           CONSTRAINT CHK_user_settings_language CHECK (language IN ('en', 'ar')) ,
    timeout INT NOT NULL,
    is_notification TINYINT DEFAULT 0,
    updated_at DATETIME,
    updated_by INT,
    CONSTRAINT FK_user_settings_users_1 FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_user_settings_timezones FOREIGN KEY (timezone_id) REFERENCES timezones(id),
    CONSTRAINT FK_user_settings_users_2 FOREIGN KEY (updated_by) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.user_settings;
-- +goose StatementEnd
