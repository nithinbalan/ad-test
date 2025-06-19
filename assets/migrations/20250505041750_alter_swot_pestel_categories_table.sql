-- +goose Up
-- +goose StatementBegin
ALTER TABLE swot_pestel_categories
ADD is_swot TINYINT NULL;

ALTER TABLE swot_pestel_categories
ADD is_pestel TINYINT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE swot_pestel_categories
DROP COLUMN is_swot;

ALTER TABLE swot_pestel_categories
DROP COLUMN is_pestel;
-- +goose StatementEnd
