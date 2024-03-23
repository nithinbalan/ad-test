-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en VARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_compliance_categories_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.compliance_categories;
-- +goose StatementEnd
