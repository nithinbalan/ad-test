-- +goose Up
-- +goose StatementBegin
CREATE TABLE governance_report_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    governance_report_id INT NOT NULL,
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
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_governance_report_documents_governance_reports FOREIGN KEY (governance_report_id) REFERENCES governance_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_governance_report_documents_user_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE governance_report_documents;
-- +goose StatementEnd
