-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_and_exercise_recovery_action_plans
ADD percentage INT NULL DEFAULT 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_and_exercise_recovery_action_plans
DROP COLUMN percentage;
-- +goose StatementEnd
