-- +goose Up
-- +goose StatementBegin
ALTER TABLE controls ADD control_automation_id INT NULL;
ALTER TABLE controls ADD control_consistency_id INT NULL;
ALTER TABLE controls ADD CONSTRAINT FK_controls_control_automations FOREIGN KEY (control_automation_id) REFERENCES control_automations(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE controls ADD CONSTRAINT FK_controls_control_consistencies FOREIGN KEY (control_consistency_id) REFERENCES control_consistencies(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE controls DROP CONSTRAINT FK_controls_control_automations;
ALTER TABLE controls DROP CONSTRAINT FK_controls_control_consistencies;
ALTER TABLE controls DROP COLUMN control_automation_id;
ALTER TABLE controls DROP COLUMN control_consistency_id;
-- +goose StatementEnd
