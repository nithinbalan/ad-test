-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_and_exercise_recovery_action_plans
ADD is_recovery TINYINT NOT NULL DEFAULT 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_and_exercise_recovery_action_plans 
    DROP COLUMN IF EXISTS is_recovery;
-- +goose StatementEnd
