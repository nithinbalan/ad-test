-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_notifications
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    user_id INT NOT NULL,
    module_id INT NOT NULL,
    icon NVARCHAR(45) NOT NULL,
    message_en NVARCHAR(255) NOT NULL,
    message_ar NVARCHAR(255) NOT NULL,
    redirect_url NVARCHAR(255) NOT NULL,
    is_new BIT DEFAULT 1,
    is_read BIT DEFAULT 0,
    created_by INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_user_notifications_users_1 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_user_notifications_modules FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_user_notifications_users_2 FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_notifications;
-- +goose StatementEnd
