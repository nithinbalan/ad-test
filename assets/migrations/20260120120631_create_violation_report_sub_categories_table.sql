-- +goose Up
-- +goose StatementBegin
CREATE TABLE violation_report_sub_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    violation_report_category_id INT NOT NULL,
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_violation_report_sub_categories_status CHECK (status IN ('active', 'inactive')),

    CONSTRAINT FK_violation_report_sub_categories_violation_report_category FOREIGN KEY (violation_report_category_id) REFERENCES violation_report_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.violation_report_sub_categories;
-- +goose StatementEnd