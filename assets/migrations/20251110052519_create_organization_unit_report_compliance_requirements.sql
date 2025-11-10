-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_unit_report_compliance_requirements (
    id INT PRIMARY KEY IDENTITY(1,1),
    organization_unit_report_id INT NOT NULL,
    compliance_requirement_title_en NVARCHAR(255) NOT NULL,
    compliance_requirement_title_ar NVARCHAR(255),
    compliance_requirement_category_title_en NVARCHAR(255),
    compliance_requirement_category_title_ar NVARCHAR(255),
    compliance_status_title_en NVARCHAR(150) NULL,
    compliance_status_title_ar NVARCHAR(150) NULL,
    compliance_status_color_code VARCHAR(9) NULL,
    severity_level_title_en NVARCHAR(100) NULL,
    severity_level_title_ar NVARCHAR(100) NULL,
    severity_level_color_code VARCHAR(9) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_organization_unit_report_compliance_requirements_reports 
        FOREIGN KEY (organization_unit_report_id) REFERENCES organization_unit_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_organization_unit_report_compliance_requirements_users_1 
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organization_unit_report_compliance_requirements;
-- +goose StatementEnd
