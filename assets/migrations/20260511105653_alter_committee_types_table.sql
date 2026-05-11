-- +goose Up
-- +goose StatementBegin

-- Drop existing check constraint
ALTER TABLE committee_types
DROP CONSTRAINT CHK_committee_types_type;

-- Recreate check constraint with new type 'other'
ALTER TABLE committee_types
ADD CONSTRAINT CHK_committee_types_type CHECK (
    type IN (
        'statutory_reference',
        'no_statutory_reference',
        'other'
    )
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Remove updated constraint
ALTER TABLE committee_types
DROP CONSTRAINT CHK_committee_types_type;

-- Restore old constraint
ALTER TABLE committee_types
ADD CONSTRAINT CHK_committee_types_type CHECK (
    type IN (
        'statutory_reference',
        'no_statutory_reference'
    )
);

-- +goose StatementEnd