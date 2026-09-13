-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_report_projects (
    id INT IDENTITY(1, 1) PRIMARY KEY,
    project_report_id INT NOT NULL,
    project_id INT NOT NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    budget DECIMAL(18, 2) NULL,
    scope NVARCHAR(MAX) NULL,
    start_date DATETIME NULL,
    end_date DATETIME NULL,
    project_type_title_en NVARCHAR(255) NULL,
    project_type_title_ar NVARCHAR(255) NULL,
    project_status_title_en NVARCHAR(255) NULL,
    project_status_title_ar NVARCHAR(255) NULL,
    project_phase_title_en NVARCHAR(255) NULL,
    project_phase_title_ar NVARCHAR(255) NULL,
    project_classification_title_en NVARCHAR(255) NULL,
    project_classification_title_ar NVARCHAR(255) NULL,
    responsible_user_name_en NVARCHAR(255) NULL,
    responsible_user_name_ar NVARCHAR(255) NULL,
    incident_action_plan_id INT NULL,
    incident_action_plan_title NVARCHAR(255) NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_report_projects_project_reports
        FOREIGN KEY (project_report_id)
        REFERENCES project_reports (id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,

    CONSTRAINT FK_project_report_projects_projects
        FOREIGN KEY (project_id)
        REFERENCES projects (id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_report_projects;
-- +goose StatementEnd
