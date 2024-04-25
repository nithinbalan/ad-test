-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description_en NVARCHAR(MAX),
    description_ar NVARCHAR(MAX),
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'not_implemented',
        CONSTRAINT CHK_compliance_statuses_type CHECK (type IN ('not_implemented', 'in_progress', 'partially_compliant', 'compliant', 'exempt', 'not_applicable')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_compliance_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.compliance_statuses;
-- +goose StatementEnd
