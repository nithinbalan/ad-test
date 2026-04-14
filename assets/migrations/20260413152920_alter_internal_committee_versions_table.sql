-- +goose Up
-- +goose StatementBegin

ALTER TABLE internal_committee_versions
ADD 
    is_disband_requested BIT NOT NULL 
        CONSTRAINT DF_internal_committee_versions_is_disband_requested DEFAULT 0;

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

-- Drop DEFAULT constraint
ALTER TABLE internal_committee_versions
DROP CONSTRAINT DF_internal_committee_versions_is_disband_requested;

-- Then drop column
ALTER TABLE internal_committee_versions
DROP COLUMN is_disband_requested;

-- +goose StatementEnd