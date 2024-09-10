-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_document_request_conversations
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    audit_document_request_id INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    commented_by_id INT NOT NULL,
    commented_at DATETIME NOT NULL,
    CONSTRAINT FK_audit_document_request_conversations_audit_document_requests FOREIGN KEY (audit_document_request_id) REFERENCES audit_document_requests(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_document_request_conversations_users_1 FOREIGN KEY (commented_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_document_request_conversations;

-- +goose StatementEnd

