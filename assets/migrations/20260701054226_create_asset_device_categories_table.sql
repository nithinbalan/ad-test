-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_device_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    is_system BIT NOT NULL DEFAULT 0,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_asset_device_categories_status
        CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE asset_device_categories;
-- +goose StatementEnd
