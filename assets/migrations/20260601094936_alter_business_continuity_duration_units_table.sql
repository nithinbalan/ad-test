-- +goose Up
-- +goose StatementBegin
-- Drop existing check constraint
ALTER TABLE business_continuity_duration_units DROP CONSTRAINT CHK_business_continuity_duration_units_type;

-- Add new check constraint with 'months' included
ALTER TABLE business_continuity_duration_units ADD CONSTRAINT CHK_business_continuity_duration_units_type CHECK (
    type IN (
        'minutes',
        'hours',
        'days',
        'weeks',
        'months'
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop check constraint with 'months'
ALTER TABLE business_continuity_duration_units DROP CONSTRAINT CHK_business_continuity_duration_units_type;

-- Add back original check constraint without 'months'
ALTER TABLE business_continuity_duration_units ADD CONSTRAINT CHK_business_continuity_duration_units_type CHECK (
    type IN (
        'minutes',
        'hours',
        'days',
        'weeks'
    )
);
-- +goose StatementEnd
