-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_workflows (
    id INT IDENTITY(1,1) PRIMARY KEY,
    project_id INT NOT NULL,
    level INT NOT NULL,
    workflow_level_type_id INT NOT NULL,
    approval_types NVARCHAR(50) NOT NULL,
    workflow_status_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_workflows_projects
        FOREIGN KEY (project_id) REFERENCES projects(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_workflows_workflow_level_types
        FOREIGN KEY (workflow_level_type_id) REFERENCES workflow_level_types(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_workflows_workflow_statuses
        FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_workflows_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_workflows;
-- +goose StatementEnd
