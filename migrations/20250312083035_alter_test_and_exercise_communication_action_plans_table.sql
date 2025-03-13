-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_and_exercise_communication_action_plans
ADD comment NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_and_exercise_communication_action_plans
DROP COLUMN comment;
-- +goose StatementEnd
