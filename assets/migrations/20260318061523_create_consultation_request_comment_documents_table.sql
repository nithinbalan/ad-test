-- +goose Up
-- +goose StatementBegin
CREATE TABLE consultation_request_comment_documents
(
    id INT PRIMARY KEY IDENTITY(1,1),
    consultation_request_comment_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,

    created_by_id INT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    -- Foreign Key Constraints
    CONSTRAINT FK_consultation_request_comment_documents
        FOREIGN KEY (consultation_request_comment_id)
        REFERENCES consultation_request_comments(id)
        ON DELETE CASCADE ON UPDATE NO ACTION,

    CONSTRAINT FK_consultation_request_comment_documents_users_1
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_consultation_request_comment_documents_users_2
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE consultation_request_comment_documents;
-- +goose StatementEnd