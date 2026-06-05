-- +goose Up
-- +goose StatementBegin
ALTER TABLE workflow_level_participants
ADD workflow_custom_role_id INT NULL;

ALTER TABLE workflow_level_participants
ADD CONSTRAINT FK_workflow_level_participants_workflow_custom_roles
FOREIGN KEY (workflow_custom_role_id)
REFERENCES workflow_custom_roles(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE workflow_level_participants
DROP CONSTRAINT FK_workflow_level_participants_workflow_custom_roles;

ALTER TABLE workflow_level_participants
DROP COLUMN workflow_custom_role_id;
-- +goose StatementEnd
