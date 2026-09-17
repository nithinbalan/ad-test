-- +goose Up
-- +goose StatementBegin
UPDATE kpi_statuses
SET color_code = '#6C757D'
WHERE type = 'new' OR title_en = 'Pending';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE kpi_statuses
SET color_code = '#808080'
WHERE type = 'new' OR title_en = 'Pending';
-- +goose StatementEnd
