-- +goose Up
-- +goose StatementBegin
CREATE TABLE integration_credentials (
    id INT PRIMARY KEY IDENTITY(1, 1),
    integration_id INT NOT NULL,
    api_key NVARCHAR(255) NULL,
    api_secret NVARCHAR(255) NULL,
    redirect_url NVARCHAR(255) NULL,
    access_token NVARCHAR(255) NULL,
    token_type NVARCHAR(255) NULL,
    scope NVARCHAR(255) NULL,
    expire_at DATETIME NULL,
    refresh_token NVARCHAR(255) NULL,
    is_active TINYINT NOT NULL DEFAULT 0,
    additional_info NVARCHAR(MAX) NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,

    CONSTRAINT FK_integration_credentials_integrations FOREIGN KEY (integration_id) REFERENCES integrations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE integration_credentials;
-- +goose StatementEnd
