-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_assessments
ADD compliance_assessment_status_id INT NULL,
    submitted_by INT NULL,
    submitted_at DATETIME NULL,
    reviewed_by INT NULL;

ALTER TABLE compliance_assessments
ADD CONSTRAINT FK_compliance_assessments_status FOREIGN KEY (compliance_assessment_status_id) REFERENCES compliance_assessment_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessments_users_4 FOREIGN KEY (submitted_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessments_users_5 FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_assessments
DROP CONSTRAINT FK_compliance_assessments_status;

ALTER TABLE compliance_assessments
DROP CONSTRAINT FK_compliance_assessments_users_4;

ALTER TABLE compliance_assessments
DROP CONSTRAINT FK_compliance_assessments_users_5;

ALTER TABLE compliance_assessments DROP COLUMN compliance_assessment_status_id;

ALTER TABLE compliance_assessments DROP COLUMN submitted_by;

ALTER TABLE compliance_assessments DROP COLUMN submitted_at;

ALTER TABLE compliance_assessments DROP COLUMN reviewed_by;
-- +goose StatementEnd
