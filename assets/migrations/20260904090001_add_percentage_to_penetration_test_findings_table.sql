-- +goose Up
-- +goose StatementBegin
ALTER TABLE penetration_test_findings
    ADD percentage INT NOT NULL CONSTRAINT DF_penetration_test_findings_percentage DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE penetration_test_findings
    DROP CONSTRAINT DF_penetration_test_findings_percentage;
ALTER TABLE penetration_test_findings
    DROP COLUMN percentage;
-- +goose StatementEnd
