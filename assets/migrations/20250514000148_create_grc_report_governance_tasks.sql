-- +goose Up
-- +goose StatementBegin
CREATE TABLE grc_report_governance_tasks (
    id INT PRIMARY KEY IDENTITY(1,1),
    grc_report_id INT NOT NULL,
    task_title NVARCHAR(255) NOT NULL,
    task_description NVARCHAR(MAX) NULL,
    responsible_user_name_en NVARCHAR(255) NOT NULL,
    responsible_user_name_ar NVARCHAR(255) NOT NULL,
    task_status NVARCHAR(100) NOT NULL,
    task_status_title_en NVARCHAR(255) NOT NULL,
    task_status_title_ar NVARCHAR(255) NOT NULL,
    task_status_color_code NVARCHAR(50) NOT NULL,
    severity_level_title_en NVARCHAR(255) NULL,
    severity_level_title_ar NVARCHAR(255) NULL,
    severity_level_color_code NVARCHAR(45) NULL,
    target_date DATE NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_grc_report_governance_tasks_grc_reports FOREIGN KEY (grc_report_id) REFERENCES grc_reports(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_grc_report_governance_tasks_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grc_report_governance_tasks;
-- +goose StatementEnd
