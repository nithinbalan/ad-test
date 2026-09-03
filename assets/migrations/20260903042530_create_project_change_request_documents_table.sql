-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_change_request_documents (
    id INT PRIMARY KEY IDENTITY(1, 1),
    project_change_request_id INT NOT NULL,
    url VARCHAR(255) NULL,
    name VARCHAR(100) NOT NULL,
    size INT NOT NULL,
    token VARCHAR(100) NOT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(100) NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_change_request_documents_change_requests FOREIGN KEY (project_change_request_id)
        REFERENCES project_change_requests(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_request_documents_users_1 FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_request_documents_users_2 FOREIGN KEY (updated_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_change_request_documents;
-- +goose StatementEnd
