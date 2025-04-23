-- +goose Up
-- +goose StatementBegin
CREATE TABLE governance_report_tasks (
    id INT PRIMARY KEY IDENTITY(1,1),
    governance_report_id INT NOT NULL,
    task_title NVARCHAR(255) NOT NULL,
    task_description NVARCHAR(MAX) NULL,
    responsible_user_name_en NVARCHAR(255) NOT NULL,
    responsible_user_name_ar NVARCHAR(255) NOT NULL,
    task_status NVARCHAR(100) NOT NULL,
    task_status_title_en NVARCHAR(255) NOT NULL,
    task_status_title_ar NVARCHAR(255) NOT NULL,
    task_status_color_code NVARCHAR(50) NOT NULL,
    target_date DATE NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_governance_report_tasks_report FOREIGN KEY (governance_report_id) REFERENCES governance_reports(id) ON DELETE CASCADE,
    CONSTRAINT FK_governance_report_tasks_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE governance_report_tasks;
-- +goose StatementEnd
