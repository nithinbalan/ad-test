-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    status VARCHAR(10) NOT NULL,
        CONSTRAINT CHK_sla_categories_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.sla_categories;
-- +goose StatementEnd
