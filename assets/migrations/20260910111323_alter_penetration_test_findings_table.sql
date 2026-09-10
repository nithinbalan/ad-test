-- +goose Up
-- +goose StatementBegin
ALTER TABLE penetration_test_findings
ALTER COLUMN description NVARCHAR(MAX) NULL;
ALTER TABLE penetration_test_findings
ALTER COLUMN recommendation NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE penetration_test_findings
ALTER COLUMN description NVARCHAR(MAX) NOT NULL;
ALTER TABLE penetration_test_findings
ALTER COLUMN recommendation NVARCHAR(MAX) NOT NULL;
-- +goose StatementEnd
