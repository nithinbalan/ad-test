-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_access_types (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'public',
        CONSTRAINT CHK_document_access_types_type CHECK (type IN ('public', 'shared')),
    color_code VARCHAR(25) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_document_access_types_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_access_types;
-- +goose StatementEnd
