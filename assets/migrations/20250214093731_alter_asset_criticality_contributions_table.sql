-- +goose Up
-- +goose StatementBegin
ALTER TABLE asset_criticality_contributions
ADD comments NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE asset_criticality_contributions
DROP COLUMN comments;
-- +goose StatementEnd
