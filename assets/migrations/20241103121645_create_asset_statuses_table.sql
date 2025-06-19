-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_asset_statuses_type CHECK (type IN ('active', 'available', 'maintenance', 'disposed')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_asset_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.asset_statuses;
-- +goose StatementEnd
