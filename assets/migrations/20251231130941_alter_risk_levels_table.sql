-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_levels 
ADD score_from INT NULL;

ALTER TABLE risk_levels 
ADD score_to INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_levels
DROP COLUMN score_from;

ALTER TABLE risk_levels
DROP COLUMN score_to;
-- +goose StatementEnd
