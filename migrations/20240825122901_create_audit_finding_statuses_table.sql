-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_finding_statuses
(
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(45) NOT NULL,
    title_ar NVARCHAR(45) NOT NULL,
    color_code VARCHAR(10) NOT NULL,
    type VARCHAR(20) NOT NULL,
    is_closed TINYINT NOT NULL DEFAULT 0,
        CONSTRAINT CHK_audit_finding_statuses_type CHECK (type IN ('new', 'in-progress', 'resolved', 'reopened', 'closed', 'deferred'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_finding_statuses;
-- +goose StatementEnd
