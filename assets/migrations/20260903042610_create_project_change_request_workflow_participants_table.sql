-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_change_request_workflow_participants (
    id INT IDENTITY(1,1) PRIMARY KEY,
    project_change_request_workflow_id INT NOT NULL,
    role_id INT NULL,
    user_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_change_request_workflow_participants_workflows
        FOREIGN KEY (project_change_request_workflow_id) REFERENCES project_change_request_workflows(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_request_workflow_participants_roles
        FOREIGN KEY (role_id) REFERENCES roles(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_request_workflow_participants_users_1
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_request_workflow_participants_users_2
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_change_request_workflow_participants;
-- +goose StatementEnd
