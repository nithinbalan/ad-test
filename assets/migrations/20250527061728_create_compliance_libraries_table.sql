-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_libraries (
    id INT PRIMARY KEY IDENTITY(1, 1),
    icon NVARCHAR(255) NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    issued_by NVARCHAR(255) NULL,
    classification NVARCHAR(255) NULL,
    document_classification NVARCHAR(255) NULL,
    disclaimer NVARCHAR(MAX) NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
);
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_libraries;
-- +goose StatementEnd