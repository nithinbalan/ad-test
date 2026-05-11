-- +goose Up
-- +goose StatementBegin
ALTER TABLE organization_units
ALTER COLUMN title_en NVARCHAR(300) NOT NULL;

ALTER TABLE organization_units
ALTER COLUMN title_ar NVARCHAR(300) NOT NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE organization_units
ALTER COLUMN title_en NVARCHAR(45) NOT NULL;

ALTER TABLE organization_units
ALTER COLUMN title_ar NVARCHAR(45) NOT NULL;
-- +goose StatementEnd