-- +goose Up
-- +goose StatementBegin
ALTER TABLE incident_action_plans
ADD delegatee_user_id INT NULL;

ALTER TABLE incident_action_plans
ADD CONSTRAINT FK_incident_action_plans_delegatee_user
FOREIGN KEY (delegatee_user_id) REFERENCES users(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE incident_action_plans
DROP CONSTRAINT FK_incident_action_plans_delegatee_user;

ALTER TABLE incident_action_plans
DROP COLUMN delegatee_user_id;
-- +goose StatementEnd
