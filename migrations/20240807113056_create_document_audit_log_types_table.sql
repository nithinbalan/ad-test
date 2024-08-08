-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_audit_log_types
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL,
    CONSTRAINT CHK_document_audit_log_types_type CHECK (type IN (
        'published', 
        'viewed', 
        'downloaded', 
        'acknowledged'
    ))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE document_audit_log_types;
-- +goose StatementEnd