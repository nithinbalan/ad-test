-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirement_responsible_units
ADD responsible_user_id INT NULL;

ALTER TABLE compliance_assessment_requirement_responsible_units
ADD CONSTRAINT FK_compliance_assessment_requirement_responsible_units_responsible_user FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirement_responsible_units
DROP CONSTRAINT FK_compliance_assessment_requirement_responsible_units_responsible_user;

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP COLUMN responsible_user_id;
-- +goose StatementEnd
