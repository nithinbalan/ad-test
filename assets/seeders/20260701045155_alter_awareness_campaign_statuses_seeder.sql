-- +goose Up
-- +goose StatementBegin
UPDATE awareness_campaign_statuses
SET color_code =
CASE type
    WHEN 'planned' THEN '#007BFF'
    WHEN 'submitted' THEN '#90DD75'
    WHEN 'cancelled' THEN '#FEC107'
    WHEN 'completed' THEN '#28A745'
    WHEN 'not_started' THEN '#6C757D'
END
WHERE type IN (
    'planned',
    'submitted',
    'cancelled',
    'completed',
    'not_started'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE awareness_campaign_statuses
SET color_code = NULL
WHERE type IN (
    'planned',
    'submitted',
    'cancelled',
    'completed',
    'not_started'
);
-- +goose StatementEnd