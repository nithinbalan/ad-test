-- +goose Up
-- +goose StatementBegin
CREATE TABLE account_types (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
        CONSTRAINT CHK_account_types_type
            CHECK (type IN ('human', 'admin', 'service', 'vendor', 'shared')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_account_types_status
            CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.account_types;
-- +goose StatementEnd
