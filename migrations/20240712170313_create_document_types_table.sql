-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_types (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en VARCHAR(45) NOT NULL,
    title_ar VARCHAR(45) NOT NULL,
    type VARCHAR(45) NOT NULL,
       CONSTRAINT CHK_document_types_type CHECK (type IN ('policy', 'procedure', 'others')),
    status VARCHAR(45) NOT NULL DEFAULT 'active',
       CONSTRAINT CHK_document_types_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_types;
-- +goose StatementEnd
