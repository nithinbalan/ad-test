-- +goose Up
-- +goose StatementBegin
-- Preserve source IDs as historical references without blocking source deletion.
ALTER TABLE penetration_test_report_documents
DROP CONSTRAINT FK_penetration_test_report_documents_supporting_documents;

ALTER TABLE penetration_test_report_documents
DROP CONSTRAINT FK_penetration_test_report_documents_retest_documents;

ALTER TABLE penetration_test_report_retest_findings
DROP CONSTRAINT FK_penetration_test_report_retest_findings_source;

ALTER TABLE penetration_test_report_retests
DROP CONSTRAINT FK_penetration_test_report_retests_source;

ALTER TABLE penetration_test_report_findings
DROP CONSTRAINT FK_penetration_test_report_findings_source;

ALTER TABLE penetration_test_report_penetration_tests
DROP CONSTRAINT FK_penetration_test_report_penetration_tests_penetration_tests;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Restoring these constraints requires the referenced source records to exist.
ALTER TABLE penetration_test_report_penetration_tests
ADD CONSTRAINT FK_penetration_test_report_penetration_tests_penetration_tests
FOREIGN KEY (penetration_test_id) REFERENCES penetration_tests (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE penetration_test_report_findings
ADD CONSTRAINT FK_penetration_test_report_findings_source
FOREIGN KEY (penetration_test_finding_id) REFERENCES penetration_test_findings (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE penetration_test_report_retests
ADD CONSTRAINT FK_penetration_test_report_retests_source
FOREIGN KEY (penetration_test_retest_id) REFERENCES penetration_test_retests (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE penetration_test_report_retest_findings
ADD CONSTRAINT FK_penetration_test_report_retest_findings_source
FOREIGN KEY (penetration_test_retest_finding_id) REFERENCES penetration_test_retest_findings (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE penetration_test_report_documents
ADD CONSTRAINT FK_penetration_test_report_documents_retest_documents
FOREIGN KEY (penetration_test_retest_finding_document_id) REFERENCES penetration_test_retest_finding_documents (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE penetration_test_report_documents
ADD CONSTRAINT FK_penetration_test_report_documents_supporting_documents
FOREIGN KEY (penetration_test_supporting_document_document_id) REFERENCES penetration_test_supporting_document_documents (id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
