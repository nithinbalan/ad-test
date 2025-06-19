-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    user_id INT NOT NULL,
    user_document_type_id INT NOT NULL,
    user_document_sub_type_id INT NOT NULL,
    issued_by NVARCHAR(MAX) NULL,
    expiry_date DATE NULL,
    remarks NVARCHAR(MAX) NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_user_documents_user_document_types FOREIGN KEY (user_document_type_id) REFERENCES user_document_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_user_documents_user_document_sub_types FOREIGN KEY (user_document_sub_type_id) REFERENCES user_document_sub_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_user_documents_users_1 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_user_documents_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_user_documents_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION


)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.user_documents;
-- +goose StatementEnd
