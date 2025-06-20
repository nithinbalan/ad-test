-- +goose Up
-- +goose StatementBegin
ALTER TABLE business_continuity_plans ADD is_latest_version TINYINT NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE business_continuity_plans DROP COLUMN is_latest_version;
-- +goose StatementEnd
