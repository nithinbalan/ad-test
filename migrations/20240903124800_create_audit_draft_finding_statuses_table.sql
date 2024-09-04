-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_draft_finding_statuses
(
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(45) NOT NULL,
    title_ar NVARCHAR(45) NOT NULL,
    color_code VARCHAR(10) NOT NULL,
    type VARCHAR(20) NOT NULL,
    is_finalized TINYINT NOT NULL DEFAULT 0,
    CONSTRAINT CHK_audit_draft_finding_statuses_type CHECK (type IN ('open', 'accepted', 'challenged', 'dismissed'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_draft_finding_statuses;
-- +goose StatementEnd
