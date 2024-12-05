-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    asset_type_id INT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_asset_categories_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_asset_categories_asset_types FOREIGN KEY (asset_type_id) REFERENCES asset_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.asset_categories;
-- +goose StatementEnd
