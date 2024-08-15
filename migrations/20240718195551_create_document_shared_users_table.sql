-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_shared_users (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_id INT NOT NULL,
    user_id INT NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'shared-user',
        CONSTRAINT CHK_document_shared_users_type CHECK (type IN ('created-user', 'review-user','responsible-user', 'shared-user')),
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_document_shared_users_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_shared_users_users_1 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_shared_users_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_shared_users;
-- +goose StatementEnd
