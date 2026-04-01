-- +goose Up
-- +goose StatementBegin

UPDATE modules
SET route = '/organization/governance-intiatives'
WHERE identifier = 'strategy_intiative';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

UPDATE modules
SET route = '/strategy/governance-intiatives'
WHERE identifier = 'strategy_intiative';

-- +goose StatementEnd