-- +goose Up
-- +goose StatementBegin
CREATE TABLE third_party_associated_documents (
    id INT PRIMARY KEY IDENTITY(1, 1),
	third_party_register_id INT NOT NULL,
    document_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
	CONSTRAINT FK_third_party_associated_documents_third_party_register FOREIGN KEY (third_party_register_id) REFERENCES third_party_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_associated_documents_associated_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_third_party_associated_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE third_party_associated_documents;

-- +goose StatementEnd
