-- +goose Up
-- +goose StatementBegin
ALTER TABLE control_action_plans
ADD review_user_id INT;

ALTER TABLE control_action_plans
ADD CONSTRAINT FK_control_action_plans_review_user FOREIGN KEY (review_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE control_action_plans
DROP CONSTRAINT FK_control_action_plans_review_user;

ALTER TABLE control_action_plans
DROP COLUMN review_user_id;
-- +goose StatementEnd
