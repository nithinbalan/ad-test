-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_and_exercise_recovery_action_plans
ADD description NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_and_exercise_recovery_action_plans
DROP COLUMN description;
-- +goose StatementEnd
