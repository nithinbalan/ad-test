-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirement_comments
ADD compliance_assessment_requirement_status_id INT NULL;

ALTER TABLE compliance_assessment_requirement_comments
ADD CONSTRAINT FK_compliance_assessment_requirement_comments_status FOREIGN KEY (compliance_assessment_requirement_status_id) REFERENCES compliance_assessment_requirement_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirement_comments
DROP CONSTRAINT FK_compliance_assessment_requirement_comments_status;

ALTER TABLE compliance_assessment_requirement_comments
DROP COLUMN compliance_assessment_requirement_status_id;
-- +goose StatementEnd
