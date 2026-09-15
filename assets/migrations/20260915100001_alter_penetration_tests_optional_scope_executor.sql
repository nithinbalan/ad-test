-- +goose Up
-- +goose StatementBegin
ALTER TABLE penetration_tests
ALTER COLUMN penetration_test_executor_id INT NULL;

ALTER TABLE penetration_tests
ALTER COLUMN scope NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE penetration_tests
ALTER COLUMN penetration_test_executor_id INT NOT NULL;

ALTER TABLE penetration_tests
ALTER COLUMN scope NVARCHAR(MAX) NOT NULL;
-- +goose StatementEnd
