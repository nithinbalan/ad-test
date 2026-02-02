-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_impact_types (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
 type VARCHAR(100) NOT NULL,
        CONSTRAINT CHK_compliance_assessment_impact_types_type CHECK (type IN ('regulatory','reputational','financial','operational')),
    description_en NVARCHAR(max) NULL,
    description_ar NVARCHAR(max) NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_compliance_assessment_impact_types_status CHECK (status IN ('active', 'inactive')),
   
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_assessment_impact_types;
-- +goose StatementEnd