-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_report_action_plans (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_report_id INT NOT NULL,
    risk_report_identification_id INT NOT NULL,

    action_plan_title NVARCHAR(255) NOT NULL,
    action_plan_description NVARCHAR(MAX) NULL,

    severity_level_title_en NVARCHAR(255) NOT NULL,
    severity_level_title_ar NVARCHAR(255) NOT NULL,

    responsible_user_name_en NVARCHAR(255) NOT NULL,
    responsible_user_name_ar NVARCHAR(255) NOT NULL,

    target_date DATE NOT NULL,

    action_plan_status_title_en NVARCHAR(255) NOT NULL,
    action_plan_status_title_ar NVARCHAR(255) NOT NULL,
    action_plan_status_color_code NVARCHAR(7) NOT NULL,

    CONSTRAINT FK_risk_report_action_plans_risk_reports FOREIGN KEY (risk_report_id) 
        REFERENCES risk_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_risk_report_action_plans_risk_report_identifications FOREIGN KEY (risk_report_identification_id) 
        REFERENCES risk_report_identifications (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_report_action_plans;
-- +goose StatementEnd
