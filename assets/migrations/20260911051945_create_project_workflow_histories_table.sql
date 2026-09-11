-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_workflow_histories (
    id INT IDENTITY(1,1) PRIMARY KEY,
    project_id INT NOT NULL,
    project_workflow_id INT NOT NULL,
    workflow_status_id INT NOT NULL,
    level INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT FK_project_workflow_histories_projects
        FOREIGN KEY (project_id) REFERENCES projects(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_workflow_histories_project_workflows
        FOREIGN KEY (project_workflow_id) REFERENCES project_workflows(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_workflow_histories_workflow_statuses
        FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_workflow_histories_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_workflow_histories;
-- +goose StatementEnd
