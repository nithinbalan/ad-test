-- +goose Up
-- +goose StatementBegin
ALTER TABLE meeting_compliance_requirements
ADD compliance_register_id INT NULL;

ALTER TABLE meeting_compliance_requirements
ADD CONSTRAINT FK_meeting_compliance_requirements_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meeting_compliance_requirements
DROP CONSTRAINT FK_meeting_compliance_requirements_compliance_register;

ALTER TABLE meeting_compliance_requirements
DROP COLUMN compliance_register_id;
-- +goose StatementEnd
