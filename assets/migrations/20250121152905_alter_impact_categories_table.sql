-- +goose Up
-- +goose StatementBegin
ALTER TABLE impact_categories 
ADD type NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE impact_categories 
DROP COLUMN type;
-- +goose StatementEnd
