-- +goose Up
-- +goose StatementBegin
ALTER TABLE kpi_register 
ADD achievement FLOAT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE kpi_register 
DROP COLUMN achievement;
-- +goose StatementEnd
