-- +goose Up

-- +goose StatementBegin

ALTER TABLE audit_statuses
DROP CONSTRAINT CHK_audit_statuses_type;

ALTER TABLE audit_statuses
ADD CONSTRAINT CHK_audit_statuses_type
CHECK (
    [type] IN (
        'completed',
        'in-progress',
        'published',
        'in-review',
        'draft',
        'planning',
        'fieldwork',
        'reporting'
    )
);

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

ALTER TABLE audit_statuses
DROP CONSTRAINT CHK_audit_statuses_type;

ALTER TABLE audit_statuses
ADD CONSTRAINT CHK_audit_statuses_type
CHECK (
    [type] IN (
        'completed',
        'in-progress',
        'published',
        'in-review',
        'draft'
    )
);

-- +goose StatementEnd