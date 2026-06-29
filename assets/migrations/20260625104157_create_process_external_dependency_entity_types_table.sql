-- +goose Up
-- +goose StatementBegin
CREATE TABLE process_external_dependency_entity_types (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    created_by_id INT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_process_external_dependency_entity_types_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_process_ext_dep_entity_types_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_process_ext_dep_entity_types_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE process_external_dependency_entity_types;
-- +goose StatementEnd
