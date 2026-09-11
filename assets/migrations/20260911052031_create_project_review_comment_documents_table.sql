-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_review_comment_documents (
    id INT IDENTITY(1,1) PRIMARY KEY,
    project_review_comment_id INT NOT NULL,
    url VARCHAR(255) NULL,
    name VARCHAR(100) NOT NULL,
    size INT NOT NULL,
    token VARCHAR(MAX) NOT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(100) NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_review_comment_documents_comments
        FOREIGN KEY (project_review_comment_id) REFERENCES project_review_comments(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_review_comment_documents_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_review_comment_documents_users_2
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_review_comment_documents;
-- +goose StatementEnd
