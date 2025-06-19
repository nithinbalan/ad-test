-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_folders (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    status VARCHAR(10) NOT NULL,
    CONSTRAINT CHK_compliance_folders_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.compliance_folders;
-- +goose StatementEnd
