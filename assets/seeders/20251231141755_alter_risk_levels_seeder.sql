-- +goose Up
-- +goose StatementBegin
UPDATE risk_levels
SET score_from = 71 , score_to = 100
WHERE type = 'low';

UPDATE risk_levels
SET score_from = 31 , score_to = 70
WHERE type = 'medium';

UPDATE risk_levels
SET score_from = 0 , score_to = 30
WHERE type = 'high';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE risk_levels
SET score_from = NULL , score_to = NULL
WHERE type = 'low';

UPDATE risk_levels
SET score_from = NULL , score_to = NULL
WHERE type = 'medium';

UPDATE risk_levels
SET score_from = NULL , score_to = NULL
WHERE type = 'high';
-- +goose StatementEnd
