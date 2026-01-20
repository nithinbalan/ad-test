-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_report_action_plans (
    id INT PRIMARY KEY IDENTITY(1,1),
    
    incident_report_id INT NOT NULL,
    incident_report_incident_id INT NOT NULL,

    action_plan_title NVARCHAR(255) NOT NULL,
    action_plan_description NVARCHAR(MAX) NULL,

    responsible_user_name_en NVARCHAR(255) NOT NULL,
    responsible_user_name_ar NVARCHAR(255) NOT NULL,

    severity_level_title_en NVARCHAR(255) NULL,
    severity_level_title_ar NVARCHAR(255) NULL,
    severity_level_color_code NVARCHAR(7) NULL,

    action_plan_status_title_en NVARCHAR(255) NOT NULL,
    action_plan_status_title_ar NVARCHAR(255) NOT NULL,
    action_plan_status_color_code NVARCHAR(7) NOT NULL,

    target_date DATETIME NOT NULL,

    CONSTRAINT FK_incident_report_action_plans_report FOREIGN KEY (incident_report_id) REFERENCES incident_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_report_action_plans_incident FOREIGN KEY (incident_report_incident_id) REFERENCES incident_report_incidents (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_report_action_plans;
-- +goose StatementEnd
