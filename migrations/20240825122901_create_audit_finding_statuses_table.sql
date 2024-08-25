-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_findings_statuses
(
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en VARCHAR(45) NOT NULL,
    title_ar VARCHAR(45) NOT NULL,
    color_code VARCHAR(10) NOT NULL,
    type VARCHAR(20) NOT NULL,
    is_closed TINYINT NOT NULL DEFAULT 0,
    CONSTRAINT CHK_audit_findings_statuses_type CHECK (type IN ('new', 'in-progress', 'resolved', 'reopened', 'closed', 'deferred'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_findings_statuses;
-- +goose StatementEnd
