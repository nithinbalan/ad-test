-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_change_request_workflow_reviews (
    id INT IDENTITY(1,1) PRIMARY KEY,
    project_change_request_workflow_id INT NOT NULL,
    workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_change_request_workflow_reviews_workflows
        FOREIGN KEY (project_change_request_workflow_id) REFERENCES project_change_request_workflows(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_request_workflow_reviews_workflow_statuses
        FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_request_workflow_reviews_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_request_workflow_reviews_users_2
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_change_request_workflow_reviews;
-- +goose StatementEnd
