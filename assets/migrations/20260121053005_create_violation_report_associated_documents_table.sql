-- +goose Up
-- +goose StatementBegin
CREATE TABLE violation_report_associated_documents (
    id INT PRIMARY KEY IDENTITY(1, 1),
    violation_report_id INT NOT NULL,
    document_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_violation_report_associated_documents_violation_report FOREIGN KEY (violation_report_id) REFERENCES violation_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_report_associated_documents_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE violation_report_associated_documents;
-- +goose StatementEnd