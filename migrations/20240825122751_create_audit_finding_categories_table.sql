-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_finding_categories
(
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(45) NOT NULL,
    title_ar NVARCHAR(45) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    type VARCHAR(20) NOT NULL,
        CONSTRAINT CHK_audit_finding_categories_type CHECK (type IN ('major-nc', 'mincor-nc', 'ofi', 'observation'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_finding_categories;
-- +goose StatementEnd
