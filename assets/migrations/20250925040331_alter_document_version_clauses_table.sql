-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_version_clauses 
ALTER COLUMN clause_title_en NVARCHAR(MAX);

ALTER TABLE document_version_clauses 
ALTER COLUMN clause_title_ar NVARCHAR(MAX);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_version_clauses 
ALTER COLUMN clause_title_en VARCHAR(255);

ALTER TABLE document_version_clauses 
ALTER COLUMN clause_title_ar NVARCHAR(255);
-- +goose StatementEnd
