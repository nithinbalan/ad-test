-- +goose Up
-- +goose StatementBegin
CREATE TABLE role_permissions (
    id INT PRIMARY KEY IDENTITY(1,1),
    role_id INT NOT NULL,
    module_id INT NOT NULL,
    role_permission_type_id INT NOT NULL,
    is_granted BIT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by INT NOT NULL,
    updated_at DATETIME NULL,
    updated_by INT NULL,
    CONSTRAINT FK_role_permissions_modules FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_role_permissions_roles FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_role_permissions_role_permission_types FOREIGN KEY (role_permission_type_id) REFERENCES role_permission_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_role_permissions_users_1 FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_role_permissions_users_2 FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE role_permissions;
-- +goose StatementEnd
