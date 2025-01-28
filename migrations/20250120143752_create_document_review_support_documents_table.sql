-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_review_support_documents
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    document_review_id INT NOT NULL,
    url VARCHAR(1000) NOT NULL,
    name NVARCHAR(255) NOT NULL,
    size INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_document_review_support_docs_document_reviews FOREIGN KEY (document_review_id) REFERENCES document_reviews(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE document_review_support_documents;
-- +goose StatementEnd
