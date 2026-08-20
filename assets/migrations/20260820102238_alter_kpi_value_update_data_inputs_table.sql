-- +goose Up
-- +goose StatementBegin
ALTER TABLE kpi_value_update_data_inputs
ALTER COLUMN value DECIMAL(18, 2) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE kpi_value_update_data_inputs SET value = 0 WHERE value IS NULL;
ALTER TABLE kpi_value_update_data_inputs ALTER COLUMN value INT NOT NULL;
-- +goose StatementEnd
