-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_compliance_requirements
ALTER COLUMN compliance_clause_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_compliance_requirements
ALTER COLUMN compliance_clause_id INT NOT NULL;
-- +goose StatementEnd
