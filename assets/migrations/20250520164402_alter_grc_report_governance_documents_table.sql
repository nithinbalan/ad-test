-- +goose Up
-- +goose StatementBegin
ALTER TABLE grc_report_governance_documents
ADD acknowledge_due_date DATE NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE grc_report_governance_documents
DROP COLUMN acknowledge_due_date;
-- +goose StatementEnd
