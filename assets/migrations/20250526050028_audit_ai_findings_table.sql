-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_ai_findings (
    id INT PRIMARY KEY IDENTITY(1,1),
    record_id NVARCHAR(100) NOT NULL,
    all_records NVARCHAR(MAX),
    all_records_status BIT DEFAULT 0,
    duplicates_records NVARCHAR(MAX),
    duplicates_records_status BIT DEFAULT 0,
    policy_violated_records NVARCHAR(MAX),
    policy_violated_status BIT DEFAULT 0,
    missing_mandatory_fields_records NVARCHAR(MAX),
    missing_mandatory_fields_status BIT DEFAULT 0,
    doubtful_transactions_records NVARCHAR(MAX),
    doubtful_transactions_status BIT DEFAULT 0,
    created_at DATETIME NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_ai_findings;
-- +goose StatementEnd
