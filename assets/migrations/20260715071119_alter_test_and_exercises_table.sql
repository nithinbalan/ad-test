-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_and_exercises
ADD is_sytem_register_updated BIT NOT NULL CONSTRAINT DF_test_and_exercises_is_sytem_register_updated DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_and_exercises
DROP CONSTRAINT DF_test_and_exercises_is_sytem_register_updated;
ALTER TABLE test_and_exercises
DROP COLUMN is_sytem_register_updated;
-- +goose StatementEnd
