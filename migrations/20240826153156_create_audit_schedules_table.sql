-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_schedules
(
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_plan_id INT NOT NULL,
    audit_plan_schedule_id INT NOT NULL,
    audit_activity_type_id INT NOT NULL,
    audit_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,
    start_date DATE NULL,
    target_date DATE NULL,
    responsible_user_id INT NULL,
    audit_schedule_status_id INT NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_audit_schedules_audit_plans FOREIGN KEY (audit_plan_id) REFERENCES audit_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_schedules_audit_plan_schedules FOREIGN KEY (audit_plan_schedule_id) REFERENCES audit_plan_schedules(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_schedules_audit_activity_types FOREIGN KEY (audit_activity_type_id) REFERENCES audit_activity_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION, 
    CONSTRAINT FK_audit_schedules_audits FOREIGN KEY (audit_id) REFERENCES audits(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_schedules_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_schedules_audit_statuses FOREIGN KEY (audit_schedule_status_id) REFERENCES audit_schedule_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_schedules_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_schedules_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_schedules;
-- +goose StatementEnd
