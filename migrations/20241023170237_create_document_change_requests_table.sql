-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_change_requests
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    documents_id INT NOT NULL,
    document_change_request_type_id INT NOT NULL,
    reason_for_change NVARCHAR(MAX) NULL,
    document_change_request_statuses_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_document_change_requests_documents FOREIGN KEY (documents_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_change_requests_types FOREIGN KEY (document_change_request_type_id) REFERENCES document_change_request_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_change_requests_statuses FOREIGN KEY (document_change_request_statuses_id) REFERENCES document_change_request_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_change_requests_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_change_requests_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE document_change_requests;
-- +goose StatementEnd
