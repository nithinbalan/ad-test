-- +goose Up
-- +goose StatementBegin
ALTER TABLE audit_ai_findings
ADD CONSTRAINT DF_audit_ai_findings_created_at
DEFAULT GETDATE() FOR created_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE audit_ai_findings
DROP CONSTRAINT DF_audit_ai_findings_created_at;
-- +goose StatementEnd
