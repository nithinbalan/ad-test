-- +goose Up
-- +goose StatementBegin
CREATE TABLE plugin_config (
    id INT PRIMARY KEY IDENTITY(1, 1),
    plugin_type NVARCHAR(100) NOT NULL,
    plugin_name NVARCHAR(255) NOT NULL,
    plugin_config NVARCHAR(MAX) NULL,
    schema_config NVARCHAR(MAX) NULL,
    plugin_docs NVARCHAR(MAX) NULL,
    samples NVARCHAR(MAX) NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    deleted_at DATETIME NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE plugin_config;
-- +goose StatementEnd
