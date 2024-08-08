-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_audit_logs
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    document_id INT NOT NULL,
    document_version_id INT NOT NULL,
    document_audit_log_type_id INT NOT NULL,
    user_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_document_audit_logs_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_audit_logs_document_versions FOREIGN KEY (document_version_id) REFERENCES document_versions(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_audit_logs_document_audit_log_types FOREIGN KEY (document_audit_log_type_id) REFERENCES document_audit_log_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_audit_logs_users_1 FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_audit_logs_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,

);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE document_audit_logs;

-- +goose StatementEnd