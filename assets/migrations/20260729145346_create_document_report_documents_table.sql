-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_report_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_report_id INT NOT NULL,
    reference_number VARCHAR(45) NOT NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    document_category_title_en NVARCHAR(255) NOT NULL,
    document_category_title_ar NVARCHAR(255) NOT NULL,
    document_type_title_en NVARCHAR(255) NOT NULL,
    document_type_title_ar NVARCHAR(255) NOT NULL,
    document_status_title_en NVARCHAR(100) NOT NULL,
    document_status_title_ar NVARCHAR(100) NOT NULL,
    document_status_color_code VARCHAR(25) NOT NULL,
    document_status_type VARCHAR(10) NOT NULL,
    document_access_type_title_en NVARCHAR(100) NOT NULL,
    document_access_type_title_ar NVARCHAR(100) NOT NULL,
    document_access_type_color_code VARCHAR(25) NOT NULL,
    confidentiality_classification_title_en NVARCHAR(100) NOT NULL,
    confidentiality_classification_title_ar NVARCHAR(100) NOT NULL,
    confidentiality_classification_type VARCHAR(150) NOT NULL,
    review_frequency_title_en NVARCHAR(100) NULL,
    review_frequency_title_ar NVARCHAR(100) NULL,
    review_frequency_type VARCHAR(20) NULL,
    responsible_unit_title_en NVARCHAR(255) NULL,
    responsible_unit_title_ar NVARCHAR(255) NULL,
    responsible_user_name_en NVARCHAR(511) NULL,
    responsible_user_name_ar NVARCHAR(511) NULL,
    version DECIMAL(4,2) NULL,
    file_name NVARCHAR(255) NULL,
    file_type VARCHAR(255) NULL,
    file_size INT NULL,
    document_created_at DATETIME NOT NULL,
    next_review_date DATE NULL,
    published_at DATETIME NULL,
    clauses_json NVARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_document_report_documents_document_reports FOREIGN KEY (document_report_id)
        REFERENCES document_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_report_documents_users_1 FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_report_documents;
-- +goose StatementEnd
