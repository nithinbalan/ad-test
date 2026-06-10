-- +goose Up
-- +goose StatementBegin

ALTER TABLE internal_committee_versions
ADD is_details_updated BIT NOT NULL
    CONSTRAINT DF_internal_committee_versions_is_details_updated DEFAULT 0;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE internal_committee_versions
DROP CONSTRAINT DF_internal_committee_versions_is_details_updated;

ALTER TABLE internal_committee_versions
DROP COLUMN is_details_updated;

-- +goose StatementEnd