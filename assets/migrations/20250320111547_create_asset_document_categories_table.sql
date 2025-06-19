-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_document_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(250) NOT NULL,
    title_ar NVARCHAR(250) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    status VARCHAR(10) NOT NULL,
        CONSTRAINT CHK_asset_document_categories_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.asset_document_categories;
-- +goose StatementEnd
