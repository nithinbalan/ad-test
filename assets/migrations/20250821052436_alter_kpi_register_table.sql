-- +goose Up
-- +goose StatementBegin
ALTER TABLE kpi_register
DROP COLUMN weightage;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE kpi_register
ADD weightage FLOAT NULL;
-- +goose StatementEnd
