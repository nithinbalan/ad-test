-- +goose Up
-- +goose StatementBegin

ALTER TABLE test_and_exercises
ADD
    is_cyber BIT NOT NULL
        CONSTRAINT DF_test_and_exercises_is_cyber DEFAULT (0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE test_and_exercises
DROP CONSTRAINT DF_test_and_exercises_is_cyber;

ALTER TABLE test_and_exercises
DROP COLUMN is_cyber;

-- +goose StatementEnd