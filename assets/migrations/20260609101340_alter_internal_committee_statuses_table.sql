-- +goose Up
-- +goose StatementBegin

ALTER TABLE internal_committee_statuses
DROP CONSTRAINT CHK_internal_committee_statuses_type;

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
        'disbanded',
        'cancelled',
        'completed'
    )
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE internal_committee_statuses
DROP CONSTRAINT CHK_internal_committee_statuses_type;

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