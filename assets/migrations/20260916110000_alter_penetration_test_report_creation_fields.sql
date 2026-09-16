-- +goose Up
-- +goose StatementBegin
UPDATE report_test
SET created_at = source_test.created_at
FROM penetration_test_report_penetration_tests AS report_test
JOIN penetration_tests AS source_test ON source_test.id = report_test.penetration_test_id
WHERE report_test.created_at IS NULL;

UPDATE report_finding
SET created_at = source_finding.created_at
FROM penetration_test_report_findings AS report_finding
JOIN penetration_test_findings AS source_finding ON source_finding.id = report_finding.penetration_test_finding_id
WHERE report_finding.created_at IS NULL;

UPDATE report_retest
SET created_at = source_retest.created_at
FROM penetration_test_report_retests AS report_retest
JOIN penetration_test_retests AS source_retest ON source_retest.id = report_retest.penetration_test_retest_id
WHERE report_retest.created_at IS NULL;

UPDATE report_document
SET created_by_id = report.created_by_id
FROM penetration_test_report_documents AS report_document
JOIN penetration_test_report_penetration_tests AS report_test ON report_test.id = report_document.penetration_test_report_penetration_test_id
JOIN penetration_test_reports AS report ON report.id = report_test.penetration_test_report_id
WHERE report_document.created_by_id IS NULL;

ALTER TABLE penetration_test_report_penetration_tests
ALTER COLUMN created_at DATETIME NOT NULL;

ALTER TABLE penetration_test_report_findings
ALTER COLUMN created_at DATETIME NOT NULL;

ALTER TABLE penetration_test_report_retests
ALTER COLUMN created_at DATETIME NOT NULL;

ALTER TABLE penetration_test_report_documents
ALTER COLUMN created_by_id INT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE penetration_test_report_documents
ALTER COLUMN created_by_id INT NULL;

ALTER TABLE penetration_test_report_retests
ALTER COLUMN created_at DATETIME NULL;

ALTER TABLE penetration_test_report_findings
ALTER COLUMN created_at DATETIME NULL;

ALTER TABLE penetration_test_report_penetration_tests
ALTER COLUMN created_at DATETIME NULL;
-- +goose StatementEnd
