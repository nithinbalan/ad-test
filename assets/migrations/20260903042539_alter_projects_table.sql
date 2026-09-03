-- +goose Up
-- +goose StatementBegin
ALTER TABLE projects ADD
    number_of_resources INT NULL,
    completion_percentage DECIMAL(5, 2) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE projects DROP COLUMN number_of_resources, completion_percentage;
-- +goose StatementEnd
