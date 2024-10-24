-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_change_request_types (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en VARCHAR(100) NOT NULL,
    title_ar VARCHAR(100) NOT NULL,
    type VARCHAR(45) NOT NULL,
       CONSTRAINT CHK_document_change_request_types CHECK (type IN ('modify', 'archive', 'restore')),
    status VARCHAR(45) NOT NULL DEFAULT 'active',
       CONSTRAINT CHK_document_change_request_types_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_change_request_types;
-- +goose StatementEnd
