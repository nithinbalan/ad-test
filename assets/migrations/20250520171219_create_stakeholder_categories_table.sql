-- +goose Up
-- +goose StatementBegin
CREATE TABLE stakeholder_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(25) NOT NULL,
        CONSTRAINT CHK_stakeholder_categories_type CHECK (type IN ('department', 'vendor','client','individual')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_stakeholder_categories_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.stakeholder_categories;
-- +goose StatementEnd

