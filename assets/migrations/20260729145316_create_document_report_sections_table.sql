-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_report_sections (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_report_id INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    order_num INT NOT NULL,

    CONSTRAINT CHK_document_report_sections_type CHECK (type IN (
        'documents',
        'clauses-controls'
    )),
    CONSTRAINT FK_document_report_sections_document_reports FOREIGN KEY (document_report_id)
        REFERENCES document_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_report_sections;
-- +goose StatementEnd
