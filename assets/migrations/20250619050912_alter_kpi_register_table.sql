-- +goose Up
-- +goose StatementBegin
ALTER TABLE kpi_register
ADD reference_number NVARCHAR(45) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE kpi_register
DROP COLUMN reference_number;
-- +goose StatementEnd
