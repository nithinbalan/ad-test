-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_securities (
    id INT PRIMARY KEY IDENTITY(1,1),
    user_id INT NOT NULL,
    password_hash VARCHAR(max) NOT NULL,
    is_mfa TINYINT DEFAULT 0,
    otp VARCHAR(45),
    otp_expiry_at DATETIME,
    updated_at DATETIME,
    CONSTRAINT FK_user_securities_users FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.user_securities;
-- +goose StatementEnd
