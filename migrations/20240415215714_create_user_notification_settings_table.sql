-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_notification_settings (
    id INT IDENTITY(1,1) PRIMARY KEY,
    user_id INT NOT NULL,
    notification_master_id INT NOT NULL,
    is_enable TINYINT NOT NULL DEFAULT 0,
    updated_at DATETIME,
    updated_by INT,
    CONSTRAINT FK_user_notification_settings_users FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_user_notification_settings_notification_masters FOREIGN KEY (notification_master_id) REFERENCES notification_masters(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_notification_settings;
-- +goose StatementEnd
