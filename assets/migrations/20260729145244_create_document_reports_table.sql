-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_reports (
    id INT PRIMARY KEY IDENTITY(1,1),
    title NVARCHAR(255) NOT NULL,
    document_category_id INT NULL,
    document_type_id INT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_document_reports_document_categories FOREIGN KEY (document_category_id)
        REFERENCES document_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_reports_document_types FOREIGN KEY (document_type_id)
        REFERENCES document_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_reports_users_1 FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_reports;
-- +goose StatementEnd
