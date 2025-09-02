-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_framework_types (
    id INT PRIMARY KEY IDENTITY(1, 1),
    type VARCHAR(50) NOT NULL CHECK (type IN ('binary', 'maturity', 'point', 'percentage')),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_compliance_assessment_framework_types_status CHECK (status IN ('active', 'inactive')),
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS compliance_assessment_framework_types;
-- +goose StatementEnd

