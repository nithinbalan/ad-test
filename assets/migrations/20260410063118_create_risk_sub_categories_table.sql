-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_sub_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),

    risk_category_id INT NOT NULL,
        CONSTRAINT FK_risk_sub_categories_risk_category_id
        FOREIGN KEY (risk_category_id)
        REFERENCES risk_categories(id)
        ON DELETE CASCADE,

    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,

    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,

    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_risk_sub_categories_status
        CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.risk_sub_categories;
-- +goose StatementEnd