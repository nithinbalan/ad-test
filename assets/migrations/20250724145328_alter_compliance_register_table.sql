-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_register ADD compliance_assessment_framework_id INT NULL;
ALTER TABLE compliance_register ADD CONSTRAINT FK_compliance_register_compliance_assessment_frameworks FOREIGN KEY (compliance_assessment_framework_id) REFERENCES compliance_assessment_frameworks(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_register DROP CONSTRAINT FK_compliance_register_compliance_assessment_frameworks;
ALTER TABLE compliance_register DROP COLUMN compliance_assessment_framework_id;
-- +goose StatementEnd
