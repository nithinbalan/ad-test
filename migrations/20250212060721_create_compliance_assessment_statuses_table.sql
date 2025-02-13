-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'draft',
        CONSTRAINT CHK_compliance_assessment_statuses_type CHECK (type IN ('draft', 'in_review', 'approved')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_compliance_assessment_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.compliance_assessment_statuses;
-- +goose StatementEnd
