-- +goose Up
-- +goose StatementBegin
CREATE TABLE violation_report_categories (
    id INT IDENTITY(1, 1) PRIMARY KEY,
    department_id INT NOT NULL,
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',

    CONSTRAINT CHK_violation_report_categories_status
        CHECK (status IN ('active', 'inactive')),

    CONSTRAINT FK_violation_report_categories_department
        FOREIGN KEY (department_id)
        REFERENCES organization_units(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE violation_report_categories;
-- +goose StatementEnd
