-- +goose Up
-- +goose StatementBegin
CREATE TABLE process_activity_users (
    id INT PRIMARY KEY IDENTITY(1,1),
    process_activity_id INT NOT NULL,
    user_id INT NOT NULL,
    responsible_unit_id INT NULL,
    role_type VARCHAR(50) NOT NULL,
    CONSTRAINT CHK_process_activity_users_role_type CHECK (
        role_type IN ('responsible', 'accountable', 'consulted', 'informed')
    ),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT FK_process_activity_users_activity FOREIGN KEY (process_activity_id) REFERENCES process_activities(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_process_activity_users_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_process_activity_users_org_unit FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS process_activity_users;
-- +goose StatementEnd
