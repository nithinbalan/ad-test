-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_change_request_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_change_request_id INT NOT NULL,
    url VARCHAR(1000) NOT NULL,
    name NVARCHAR(255) NOT NULL,
    size INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_document_change_request_documents_document_change_request FOREIGN KEY (document_change_request_id) REFERENCES document_change_requests(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_change_request_documents;
-- +goose StatementEnd
