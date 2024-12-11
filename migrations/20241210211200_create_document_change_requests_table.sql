-- +goose Up
-- +goose StatementBegin
CREATE TABLE documents_change_requests (
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(45) NOT NULL,
    document_id INT NOT NULL,
    document_change_request_initiate_id INT NULL,
    document_change_request_type_id INT NOT NULL,
    document_change_request_status_id INT NOT NULL,
    reason_for_change NVARCHAR(MAX) NULL,
    version DECIMAL(4, 2) NOT NULL,
    submitted_by_id INT NULL,
    submitted_at DATETIME NULL,
    published_by_id INT NULL,
    published_at DATETIME NULL,
    next_review_level INT NULL,
    is_locked TINYINT NOT NULL DEFAULT 0,
    last_checkout_at DATE NULL,
    last_checkout_by_id INT NULL,
    last_checkin_at DATE NULL,
    last_checkin_by_id INT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    -- Constraints
    CONSTRAINT FK_documents_change_requests_document FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_change_requests_document_change_request_initiate FOREIGN KEY (document_change_request_initiate_id) REFERENCES document_change_request_initiates(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_change_requests_documents_change_request_type FOREIGN KEY (document_change_request_type_id) REFERENCES document_change_request_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_change_requests_documents_change_request_status FOREIGN KEY (document_change_request_status_id) REFERENCES document_change_request_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_change_requests_user_1 FOREIGN KEY (submitted_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_change_requests_user_2 FOREIGN KEY (published_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_change_requests_user_3 FOREIGN KEY (last_checkout_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_change_requests_user_4 FOREIGN KEY (last_checkin_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_change_requests_user_5 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_documents_change_requests_user_6 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE documents_change_requests;
-- +goose StatementEnd
