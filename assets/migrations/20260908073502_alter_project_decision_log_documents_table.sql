-- +goose Up
-- +goose StatementBegin
ALTER TABLE project_decision_log_documents 
ALTER COLUMN url VARCHAR(1000) NULL;

ALTER TABLE project_decision_log_documents 
ALTER COLUMN name NVARCHAR(255) NOT NULL;

ALTER TABLE project_decision_log_documents 
ALTER COLUMN type VARCHAR(255) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_decision_log_documents 
ALTER COLUMN url VARCHAR(255) NULL;

ALTER TABLE project_decision_log_documents 
ALTER COLUMN name VARCHAR(100) NOT NULL;

ALTER TABLE project_decision_log_documents 
ALTER COLUMN type VARCHAR(100) NOT NULL;
-- +goose StatementEnd

