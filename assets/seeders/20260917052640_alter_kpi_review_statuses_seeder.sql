-- +goose Up
-- +goose StatementBegin
UPDATE kpi_review_statuses
SET color_code = '#90DD75'
WHERE type = 'updated';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE kpi_review_statuses
SET color_code = '#00AF50'
WHERE type = 'updated';
-- +goose StatementEnd
