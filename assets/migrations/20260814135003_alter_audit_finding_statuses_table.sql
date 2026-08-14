-- +goose Up
-- +goose StatementBegin

ALTER TABLE audit_finding_statuses
DROP CONSTRAINT CHK_audit_finding_statuses_type;

ALTER TABLE audit_finding_statuses
ALTER COLUMN type VARCHAR(30) NOT NULL;

ALTER TABLE audit_finding_statuses
ADD CONSTRAINT CHK_audit_finding_statuses_type
CHECK (
    [type] IN (
        'deferred',
        'closed',
        'reopened',
        'resolved',
        'in-progress',
        'new',
        're-assessment-pending'
    )
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE audit_finding_statuses
DROP CONSTRAINT CHK_audit_finding_statuses_type;

DELETE FROM audit_finding_statuses
WHERE [type] = 're-assessment-pending';

ALTER TABLE audit_finding_statuses
ALTER COLUMN type VARCHAR(20) NOT NULL;

ALTER TABLE audit_finding_statuses
ADD CONSTRAINT CHK_audit_finding_statuses_type
CHECK (
    [type] IN (
        'deferred',
        'closed',
        'reopened',
        'resolved',
        'in-progress',
        'new'
    )
);

-- +goose StatementEnd