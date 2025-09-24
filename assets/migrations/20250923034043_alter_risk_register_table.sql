-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_register
ADD next_review_level INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_register
DROP COLUMN next_review_level;
-- +goose StatementEnd
