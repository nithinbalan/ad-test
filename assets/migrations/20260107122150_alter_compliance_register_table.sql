-- +goose Up
-- +goose StatementBegin

ALTER TABLE compliance_register
ALTER COLUMN compliance_status_id INT NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_register
ALTER COLUMN compliance_status_id INT NULL;

-- +goose StatementEnd
