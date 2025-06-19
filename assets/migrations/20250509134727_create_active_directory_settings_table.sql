-- +goose Up
-- +goose StatementBegin
CREATE TABLE active_directory_settings (
    id INT PRIMARY KEY IDENTITY(1,1),
    server VARCHAR(255) NOT NULL,
    port INT NOT NULL,
    use_ssl BIT NOT NULL DEFAULT 0,
    base_dn NVARCHAR(500) NOT NULL,
    bind_dn NVARCHAR(500) NOT NULL,
    bind_password VARCHAR(max) NOT NULL,
    user_search_attr VARCHAR(100) NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(max) NOT NULL,
    created_by_id INT,
    updated_by_id INT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,

    CONSTRAINT FK_active_directory_settings_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_active_directory_settings_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE active_directory_settings;
-- +goose StatementEnd
