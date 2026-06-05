-- +goose Up
-- +goose StatementBegin

ALTER TABLE internal_committee_workflow_participants
ADD custom_role_id INT NULL;

ALTER TABLE internal_committee_workflow_participants
ADD CONSTRAINT FK_internal_committee_workflow_participants_custom_role
FOREIGN KEY (custom_role_id)
REFERENCES workflow_custom_roles(id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE internal_committee_workflow_participants
DROP CONSTRAINT FK_internal_committee_workflow_participants_custom_role;

ALTER TABLE internal_committee_workflow_participants
DROP COLUMN custom_role_id;

-- +goose StatementEnd