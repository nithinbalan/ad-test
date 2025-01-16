-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_roles (
    id INT PRIMARY KEY IDENTITY(1,1),
    business_continuity_plan_id INT NOT NULL,
    responsible_user_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    responsibility NVARCHAR(MAX) NOT NULL,
    backup_user_id INT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL DEFAULT GETDATE(),
    CONSTRAINT FK_bcp_roles_business_continuity_plans FOREIGN KEY (business_continuity_plan_id) REFERENCES business_continuity_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_roles_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_roles_users_2 FOREIGN KEY (backup_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_roles_users_3 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_roles;
-- +goose StatementEnd
