-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_associated_documents (
    id INT PRIMARY KEY IDENTITY(1, 1),
    document_id INT NOT NULL,
	associated_document_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
  
    CONSTRAINT FK_document_associated_documents_documents_2 FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_document_associated_documents_documents_1 FOREIGN KEY (associated_document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_associated_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE document_associated_documents;
-- +goose StatementEnd