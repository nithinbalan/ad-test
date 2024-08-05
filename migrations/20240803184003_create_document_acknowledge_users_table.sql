-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_acknowledge_users
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    document_id INT NOT NULL,
    user_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_document_acknowledge_users_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_acknowledge_users_users_1 FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_acknowledge_users_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,

);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE document_acknowledge_users;

-- +goose StatementEnd