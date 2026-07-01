-- +goose Up
-- +goose StatementBegin
ALTER TABLE awareness_campaign_statuses
ADD color_code VARCHAR(25) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE awareness_campaign_statuses
DROP COLUMN color_code;
-- +goose StatementEnd
