-- +goose Up
-- +goose StatementBegin
UPDATE awareness_campaign_statuses
SET
    title_en = 'Completed',
    type = 'completed'
WHERE type = 'complete';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE awareness_campaign_statuses
SET
    title_en = 'Complete',
    type = 'complete'
WHERE type = 'completed';
-- +goose StatementEnd
