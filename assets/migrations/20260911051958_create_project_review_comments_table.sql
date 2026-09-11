-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_review_comments (
    id INT IDENTITY(1,1) PRIMARY KEY,
    comment NVARCHAR(MAX) NOT NULL,
    is_edited BIT DEFAULT 0 NOT NULL,
    comment_status_id INT NOT NULL,
    project_id INT NOT NULL,
    level INT NULL,
    parent_comment_id INT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_review_comments_comment_statuses
        FOREIGN KEY (comment_status_id) REFERENCES comment_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_review_comments_projects
        FOREIGN KEY (project_id) REFERENCES projects(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_review_comments_parent_comments
        FOREIGN KEY (parent_comment_id) REFERENCES project_review_comments(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_review_comments_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_review_comments_users_2
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_review_comments;
-- +goose StatementEnd
