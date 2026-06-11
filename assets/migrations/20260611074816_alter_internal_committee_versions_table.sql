-- +goose Up
-- +goose StatementBegin

ALTER TABLE internal_committee_versions
ADD term_start_date DATE NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE internal_committee_versions
DROP COLUMN term_start_date;

-- +goose StatementEnd