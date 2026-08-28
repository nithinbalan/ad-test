-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_and_exercises ADD objective NVARCHAR(MAX) NULL;
ALTER TABLE test_and_exercises ADD scenario NVARCHAR(MAX) NULL;
ALTER TABLE test_and_exercises ADD success_criteria NVARCHAR(MAX) NULL;
-- +goose StatementEnd  

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_and_exercises DROP COLUMN success_criteria;
ALTER TABLE test_and_exercises DROP COLUMN scenario;
ALTER TABLE test_and_exercises DROP COLUMN objective;
-- +goose StatementEnd
