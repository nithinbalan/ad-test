-- +goose Up
-- +goose StatementBegin

UPDATE committee_terms
SET color_code = '#FFC107'  -- Orange
WHERE type = 'temporary';

UPDATE committee_terms
SET color_code = '#007BFF'  -- Blue
WHERE type = 'permanent';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

UPDATE committee_terms
SET color_code = NULL
WHERE type IN ('temporary', 'permanent');

-- +goose StatementEnd