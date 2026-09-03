-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_change_request_review_comment_documents (
    id INT IDENTITY(1,1) PRIMARY KEY,
    project_change_request_review_comment_id INT NOT NULL,
    name NVARCHAR(255) NOT NULL,
    token NVARCHAR(MAX) NOT NULL,
    original_name NVARCHAR(255) NOT NULL,
    type NVARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT FK_project_change_request_review_comment_documents_comments
        FOREIGN KEY (project_change_request_review_comment_id) REFERENCES project_change_request_review_comments(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_request_review_comment_documents_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_change_request_review_comment_documents;
-- +goose StatementEnd
