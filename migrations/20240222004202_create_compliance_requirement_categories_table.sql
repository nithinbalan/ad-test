-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_requirement_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    status NVARCHAR(10) NOT NULL,
        CONSTRAINT CHK_compliance_requirement_categories CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.compliance_requirement_categories;
-- +goose StatementEnd
