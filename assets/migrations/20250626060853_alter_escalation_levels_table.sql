-- +goose Up
-- +goose StatementBegin
ALTER TABLE escalation_levels
ADD designation_id INT NULL;

ALTER TABLE escalation_levels
ADD CONSTRAINT FK_escalation_levels_designations
FOREIGN KEY (designation_id) REFERENCES designations(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE escalation_levels
DROP CONSTRAINT FK_escalation_levels_designations;

ALTER TABLE escalation_levels
DROP COLUMN designation_id;
-- +goose StatementEnd
