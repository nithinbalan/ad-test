-- +goose Up
-- +goose StatementBegin

ALTER TABLE workflow_level_participants
ADD unit_id INT NULL;

ALTER TABLE workflow_level_participants
ADD CONSTRAINT FK_workflow_level_participants_organization_units
FOREIGN KEY (unit_id)
REFERENCES organization_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE workflow_level_participants
DROP CONSTRAINT FK_workflow_level_participants_organization_units;

ALTER TABLE workflow_level_participants
DROP COLUMN unit_id;

-- +goose StatementEnd