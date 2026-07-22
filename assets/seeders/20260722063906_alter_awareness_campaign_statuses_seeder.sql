-- +goose Up
-- +goose StatementBegin
UPDATE awareness_campaign_statuses
SET color_code = '#FFC107'
WHERE type = 'cancelled';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE awareness_campaign_statuses
SET color_code = '#FEC107'
WHERE type = 'cancelled';
-- +goose StatementEnd
