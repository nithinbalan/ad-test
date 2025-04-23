-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_sub_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    asset_category_id INT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_asset_sub_categories_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_asset_sub_categories_asset_categories FOREIGN KEY (asset_category_id) REFERENCES asset_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.asset_sub_categories;
-- +goose StatementEnd
