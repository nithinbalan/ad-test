-- +goose Up
-- +goose StatementBegin
CREATE TABLE processes (
    id INT PRIMARY KEY IDENTITY(1,1),
    process_group_id INT NOT NULL,
    owner_unit_id INT NOT NULL,
    owner_id INT NOT NULL,
    lifecycle_status_id INT NOT NULL,
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX),
    marked_as_auditable_item BIT NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    updated_by_id INT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT FK_processes_process_groups FOREIGN KEY (process_group_id) REFERENCES process_groups(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_processes_organization_units FOREIGN KEY (owner_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_processes_users_1 FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_processes_lifecycle_statuses FOREIGN KEY (lifecycle_status_id) REFERENCES lifecycle_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_processes_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_processes_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.processes;
-- +goose StatementEnd
