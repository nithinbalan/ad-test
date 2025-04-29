-- +goose Up
-- +goose StatementBegin
CREATE TABLE grc_report_governance_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    grc_report_id INT NOT NULL,
    reference_number NVARCHAR(50) NOT NULL,
    document_title_en NVARCHAR(255) NOT NULL,
    document_title_ar NVARCHAR(255) NOT NULL,
    responsible_user_name_en NVARCHAR(255) NOT NULL,
    responsible_user_name_ar NVARCHAR(255) NOT NULL,
    responsible_unit_title_en NVARCHAR(255) NOT NULL,
    responsible_unit_title_ar NVARCHAR(255) NOT NULL,
    total_employees INT NOT NULL DEFAULT 0,
    acknowledged_count INT NOT NULL DEFAULT 0,
    pending_count INT NOT NULL DEFAULT 0,
    document_status_title_en NVARCHAR(255) NOT NULL,
    document_status_title_ar NVARCHAR(255) NOT NULL,
    document_status_color_code NVARCHAR(50) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_grc_report_governance_documents_grc_reports FOREIGN KEY (grc_report_id) REFERENCES grc_reports(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_grc_report_governance_documents_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grc_report_governance_documents;
-- +goose StatementEnd
