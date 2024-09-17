-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_schedule_status_updates (
    id INT PRIMARY KEY IDENTITY(1, 1),
    audit_schedule_id INT NOT NULL,
    audit_schedule_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    percentage INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_audit_schedule_status_updates_audit_schedules FOREIGN KEY (audit_schedule_id) REFERENCES audit_schedules (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_schedule_status_updates_audit_schedule_statuses FOREIGN KEY (aaudit_schedule_status_id) REFERENCES audit_schedule_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_schedule_status_updates_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_schedule_status_updates_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_schedule_status_updates;
-- +goose StatementEnd
