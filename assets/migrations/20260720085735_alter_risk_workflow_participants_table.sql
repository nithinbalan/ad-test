-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_workflow_participants
ADD workflow_status_id INT NULL;

ALTER TABLE risk_workflow_participants
ADD CONSTRAINT FK_risk_workflow_participants_workflow_status
FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_workflow_participants
DROP CONSTRAINT FK_risk_workflow_participants_workflow_status;

ALTER TABLE risk_workflow_participants
DROP COLUMN workflow_status_id;
-- +goose StatementEnd
