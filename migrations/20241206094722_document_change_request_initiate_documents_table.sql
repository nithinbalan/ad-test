-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_change_request_initiate_documents
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    document_change_request_initiate_id INT NOT NULL,
    url VARCHAR(1000) NOT NULL,
    name NVARCHAR(255) NOT NULL,
    size INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_document_change_request_initiate_docs_request FOREIGN KEY (document_change_request_initiate_id) REFERENCES document_change_request_initiates(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_change_request_initiate_docs_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_change_request_initiate_docs_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE document_change_request_initiate_documents;
-- +goose StatementEnd
