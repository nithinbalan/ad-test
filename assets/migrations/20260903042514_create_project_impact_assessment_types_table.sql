-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_impact_assessment_types (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
        CONSTRAINT CHK_project_impact_assessment_types_type CHECK (type IN ('timeline', 'budget', 'resources', 'technical_complexity')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_project_impact_assessment_types_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_impact_assessment_types;
-- +goose StatementEnd
