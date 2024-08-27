-- +goose Up
-- +goose StatementBegin
ALTER TABLE audit_findings
ADD 
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE audit_findings
DROP COLUMN 
    reference_number,
    title,
    description;
-- +goose StatementEnd
