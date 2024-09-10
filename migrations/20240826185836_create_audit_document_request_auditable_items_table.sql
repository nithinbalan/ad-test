-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_document_request_auditable_items
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    audit_document_request_id INT NOT NULL,
    auditable_item_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_audit_document_request_auditable_items_audit_document_requests FOREIGN KEY (audit_document_request_id) REFERENCES audit_document_requests(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_document_request_auditable_items_auditable_items FOREIGN KEY (auditable_item_id) REFERENCES auditable_items(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_document_request_auditable_items_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,

);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_document_request_auditable_items;

-- +goose StatementEnd