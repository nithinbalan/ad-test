-- +goose Up
-- +goose StatementBegin
ALTER TABLE control_tests ADD control_automation_id INT NULL;
ALTER TABLE control_tests ADD control_consistency_id INT NULL;
ALTER TABLE control_tests ADD CONSTRAINT FK_control_tests_control_automations FOREIGN KEY (control_automation_id) REFERENCES control_automations(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE control_tests ADD CONSTRAINT FK_control_tests_control_consistencies FOREIGN KEY (control_consistency_id) REFERENCES control_consistencies(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE control_tests DROP CONSTRAINT FK_control_tests_control_automations;
ALTER TABLE control_tests DROP CONSTRAINT FK_control_tests_control_consistencies;
ALTER TABLE control_tests DROP COLUMN control_automation_id;
ALTER TABLE control_tests DROP COLUMN control_consistency_id;
-- +goose StatementEnd
