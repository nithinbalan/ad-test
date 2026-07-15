-- +goose Up
-- +goose StatementBegin

ALTER TABLE bia_process_workflow_participants
ADD custom_role_id INT NULL;

ALTER TABLE bia_process_workflow_participants
ADD CONSTRAINT FK_bia_process_workflow_participants_custom_role
FOREIGN KEY (custom_role_id)
REFERENCES workflow_custom_roles(id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE bia_process_workflow_participants
DROP CONSTRAINT FK_bia_process_workflow_participants_custom_role;

ALTER TABLE bia_process_workflow_participants
DROP COLUMN custom_role_id;

-- +goose StatementEnd