-- +goose Up
-- +goose StatementBegin

-- Drop old check constraint
ALTER TABLE internal_committee_statuses
DROP CONSTRAINT CHK_internal_committee_statuses_type;

-- Add new check constraint with under_nomination included
ALTER TABLE internal_committee_statuses
ADD CONSTRAINT CHK_internal_committee_statuses_type CHECK (
    type IN (
        'draft',
        'submitted',
        'under_nomination',
        'in_review',
        'in_progress',
        'approved',
        'revert',
        'reformation',
        'disbanded'
    )
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Drop new constraint
ALTER TABLE internal_committee_statuses
DROP CONSTRAINT CHK_internal_committee_statuses_type;

-- Restore old constraint (without under_nomination)
ALTER TABLE internal_committee_statuses
ADD CONSTRAINT CHK_internal_committee_statuses_type CHECK (
    type IN (
        'draft',
        'submitted',
        'in_review',
        'in_progress',
        'approved',
        'revert',
        'reformation',
        'disbanded'
    )
);

-- +goose StatementEnd