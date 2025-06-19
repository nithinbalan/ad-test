-- +goose Up
-- +goose StatementBegin
ALTER TABLE strategy_profile_settings 
ADD type VARCHAR(25) NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE strategy_profile_settings 
DROP COLUMN type;
-- +goose StatementEnd
