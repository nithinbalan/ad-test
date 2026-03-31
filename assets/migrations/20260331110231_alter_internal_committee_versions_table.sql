-- +goose Up
-- +goose StatementBegin

ALTER TABLE internal_committee_versions
ADD committee_authorities NVARCHAR(MAX) NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE internal_committee_versions
DROP COLUMN committee_authorities;

-- +goose StatementEnd
