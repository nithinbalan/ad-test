-- +goose Up
-- +goose StatementBegin
CREATE TABLE data_sources (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
        CONSTRAINT CHK_data_sources_type CHECK (type IN ('user_input', 'internal_module', 'external_api','ai_generated','file_upload')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_data_sources_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.data_sources;
-- +goose StatementEnd
