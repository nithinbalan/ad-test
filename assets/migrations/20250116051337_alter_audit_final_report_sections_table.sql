-- +goose Up
-- +goose StatementBegin
ALTER TABLE audit_final_report_sections
DROP CONSTRAINT CHK_audit_final_report_sections_type;

ALTER TABLE audit_final_report_sections
ADD CONSTRAINT CHK_audit_final_report_sections_type 
CHECK (type IN ('overview', 'finding-summary', 'findings', 'conclusion', 'custom', 'appendix'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE audit_final_report_sections
DROP CONSTRAINT CHK_audit_final_report_sections_type;

ALTER TABLE audit_final_report_sections
ADD CONSTRAINT CHK_audit_final_report_sections_type 
CHECK (type IN ('overview', 'finding-summary', 'findings', 'conclusion', 'custom'));
-- +goose StatementEnd
