-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS project_security_gate_checklist_exception_workflow_review_documents;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- this table no longer needed
-- +goose StatementEnd
