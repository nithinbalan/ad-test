-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_login_history (
    id INT PRIMARY KEY IDENTITY(1,1),
    user_id INT NOT NULL,
    login_at DATETIME NOT NULL,
    login_status VARCHAR(20) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    location VARCHAR(255),
    device_info VARCHAR(255),
    browser_info VARCHAR(255),
    CONSTRAINT FK_user_logins_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.user_login_history;
-- +goose StatementEnd
