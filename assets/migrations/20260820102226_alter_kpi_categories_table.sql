-- +goose Up
-- +goose StatementBegin
ALTER TABLE kpi_categories
ADD type NVARCHAR(50) NOT NULL
    CONSTRAINT DF_kpi_categories_type DEFAULT 'other';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE kpi_categories DROP CONSTRAINT DF_kpi_categories_type;
ALTER TABLE kpi_categories DROP COLUMN type;
-- +goose StatementEnd
