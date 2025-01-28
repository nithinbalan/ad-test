-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_and_exercise_statuses
DROP CONSTRAINT CHK_test_and_exercise_statuses_type;

ALTER TABLE test_and_exercise_statuses
ADD CONSTRAINT CHK_test_and_exercise_statuses_type CHECK (type IN ('scheduled', 'completed','in-progress'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_and_exercise_statuses
DROP CONSTRAINT CHK_test_and_exercise_statuses_type;

ALTER TABLE test_and_exercise_statuses
ADD CONSTRAINT CHK_test_and_exercise_statuses_type CHECK (type IN ('scheduled', 'completed'));
-- +goose StatementEnd
