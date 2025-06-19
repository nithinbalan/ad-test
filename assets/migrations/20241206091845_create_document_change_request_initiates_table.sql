-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_change_request_initiates
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    document_id INT NOT NULL,
    document_change_request_initiate_type_id INT NOT NULL,
    reason_for_change NVARCHAR(MAX) NULL,
    document_change_request_initiate_status_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_document_change_request_initiates_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_change_request_initiates_types FOREIGN KEY (document_change_request_initiate_type_id) REFERENCES document_change_request_initiate_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_change_request_initiates_statuses FOREIGN KEY (document_change_request_initiate_status_id) REFERENCES document_change_request_initiate_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_change_request_initiates_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_change_request_initiates_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE document_change_request_initiates;
-- +goose StatementEnd
