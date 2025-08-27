-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_controls
ADD completed_by_user_id INT NULL,
    completed_at DATETIME NULL,
    compliance_assessment_framework_value_id INT NULL,
    comments NVARCHAR(MAX) NULL;

ALTER TABLE compliance_controls
ADD CONSTRAINT FK_compliance_controls_completed_by_user 
    FOREIGN KEY (completed_by_user_id) REFERENCES users(id) 
    ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE compliance_controls
ADD CONSTRAINT FK_compliance_controls_framework_value 
    FOREIGN KEY (compliance_assessment_framework_value_id) REFERENCES compliance_assessment_framework_values(id) 
    ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_controls
DROP CONSTRAINT FK_compliance_controls_completed_by_user;

ALTER TABLE compliance_controls
DROP CONSTRAINT FK_compliance_controls_framework_value;

ALTER TABLE compliance_controls
DROP COLUMN completed_by_user_id,
            completed_at,
            compliance_assessment_framework_value_id,
            comments;
-- +goose StatementEnd
