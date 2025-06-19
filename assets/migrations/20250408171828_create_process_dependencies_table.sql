-- +goose Up
-- +goose StatementBegin
CREATE TABLE process_dependencies (
    id INT PRIMARY KEY IDENTITY(1, 1),
    process_id INT NOT NULL,
    process_dependency_sub_type_id INT NOT NULL,
    dependent_process_id INT NULL,
    organization_unit_id INT NULL,
    asset_id INT NULL,
    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_process_dependencies_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_process_dependencies_sub_types FOREIGN KEY (process_dependency_sub_type_id) REFERENCES process_dependency_sub_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_process_dependencies_dependent_processes FOREIGN KEY (dependent_process_id) REFERENCES processes(id) ON DELETE SET NULL ON UPDATE NO ACTION,
    CONSTRAINT FK_process_dependencies_organization_units FOREIGN KEY (organization_unit_id) REFERENCES organization_units(id) ON DELETE SET NULL ON UPDATE NO ACTION,
    CONSTRAINT FK_process_dependencies_asset_register FOREIGN KEY (asset_id) REFERENCES asset_register(id) ON DELETE SET NULL ON UPDATE NO ACTION,
    CONSTRAINT FK_process_dependencies_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_process_dependencies_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE process_dependencies;
-- +goose StatementEnd
