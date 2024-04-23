-- +goose Up
-- +goose StatementBegin
CREATE TABLE activity_logs (
    id INT PRIMARY KEY IDENTITY(1, 1),
    action VARCHAR(50) NOT NULL,
    table_name VARCHAR(45) NOT NULL,
    record_id INT,
    changed_data NVARCHAR(MAX),
    user_id INT,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
        CONSTRAINT CHK_activity_logs_action CHECK (action IN (
            'view', 'download', 'import', 'export', 'create', 'update', 'delete', 'activate', 'deactivate',
            'login', 'logout', 'access', 'edit_permissions', 'change_role', 'approve', 'reject',
            'submit', 'schedule', 'generate_report', 'upload', 'remove_file',
            'change_setting', 'restore', 'archive', 'unarchive', 'assign_task', 'complete_task', 'escalate',
            'comment', 'update_note', 'send_notification', 'lock', 'unlock'
        ))
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.activity_logs;
-- +goose StatementEnd
