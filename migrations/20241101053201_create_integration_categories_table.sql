-- +goose Up
-- +goose StatementBegin
CREATE TABLE integration_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description_en NVARCHAR(MAX) NOT NULL,
    description_ar NVARCHAR(MAX) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_integration_categories_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS integration_categories;
-- +goose StatementEnd
