-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_reviews (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_id INT NOT NULL,
    review_date DATE NOT NULL,
    is_change_required TINYINT NOT NULL DEFAULT 0,
    description  NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_document_reviews_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_document_reviews_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_reviews_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE document_reviews;
-- +goose StatementEnd
