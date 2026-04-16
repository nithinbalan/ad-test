-- +goose Up
-- +goose StatementBegin

ALTER TABLE control_tests
ALTER COLUMN control_efficiency_status_id INT NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE control_tests
ALTER COLUMN control_efficiency_status_id INT NULL;

-- +goose StatementEnd
