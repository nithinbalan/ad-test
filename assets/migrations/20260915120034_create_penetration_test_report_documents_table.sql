-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_report_documents (
    id INT PRIMARY KEY IDENTITY(1,1),

    penetration_test_report_penetration_test_id INT NOT NULL,
    penetration_test_supporting_document_document_id INT NULL,
    penetration_test_retest_finding_document_id INT NULL,
    penetration_test_report_retest_finding_id INT NULL,

    url VARCHAR(1000) NULL,
    name NVARCHAR(255) NULL,
    size INT NULL,
    token VARCHAR(255) NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) NULL,

    created_by_id INT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_penetration_test_report_documents_penetration_tests FOREIGN KEY (penetration_test_report_penetration_test_id) REFERENCES penetration_test_report_penetration_tests (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_report_documents_supporting_documents FOREIGN KEY (penetration_test_supporting_document_document_id) REFERENCES penetration_test_supporting_document_documents (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_report_documents_retest_documents FOREIGN KEY (penetration_test_retest_finding_document_id) REFERENCES penetration_test_retest_finding_documents (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_report_documents_retest_findings FOREIGN KEY (penetration_test_report_retest_finding_id) REFERENCES penetration_test_report_retest_findings (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_report_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_report_documents_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE penetration_test_report_documents;
-- +goose StatementEnd
