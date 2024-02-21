-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(45) NOT NULL,
    title_ar NVARCHAR(45) NOT NULL,
    status NVARCHAR(10) NOT NULL,
        CONSTRAINT CHK_compliance_categories_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.compliance_categories;
-- +goose StatementEnd
