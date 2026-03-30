-- +goose Up
-- +goose StatementBegin

ALTER TABLE strategic_intiatives
ADD objectives_en NVARCHAR(MAX) NULL,
    objectives_ar NVARCHAR(MAX) NULL,
    recommendations_en NVARCHAR(MAX) NULL,
    recommendations_ar NVARCHAR(MAX) NULL,
    owner_department_id INT NULL,
    strategic_initiative_status_id INT NULL,
    completion_percentage DECIMAL(5,2) NULL,
    kpis NVARCHAR(MAX) NULL,
    expected_benefits_en NVARCHAR(MAX) NULL,
    expected_benefits_ar NVARCHAR(MAX) NULL,
    budget DECIMAL(18,2) NULL,
    resource_requirements_en NVARCHAR(MAX) NULL,
    resource_requirements_ar NVARCHAR(MAX) NULL,
    dependencies_en NVARCHAR(MAX) NULL,
    dependencies_ar NVARCHAR(MAX) NULL;

ALTER TABLE strategic_intiatives
ADD CONSTRAINT FK_strategic_intiatives_owner_department
    FOREIGN KEY (owner_department_id) REFERENCES organization_units(id);

ALTER TABLE strategic_intiatives
ADD CONSTRAINT FK_strategic_intiatives_status
    FOREIGN KEY (strategic_initiative_status_id) REFERENCES strategic_initiative_statuses(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE strategic_intiatives DROP CONSTRAINT FK_strategic_intiatives_status;
ALTER TABLE strategic_intiatives DROP CONSTRAINT FK_strategic_intiatives_owner_department;

ALTER TABLE strategic_intiatives
DROP COLUMN objectives_en,
            objectives_ar,
            recommendations_en,
            recommendations_ar,
            owner_department_id,
            strategic_initiative_status_id,
            completion_percentage,
            kpis,
            expected_benefits_en,
            expected_benefits_ar,
            budget,
            resource_requirements_en,
            resource_requirements_ar,
            dependencies_en,
            dependencies_ar;

-- +goose StatementEnd