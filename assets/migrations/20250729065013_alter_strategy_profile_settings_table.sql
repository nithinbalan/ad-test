-- +goose Up
-- +goose StatementBegin
ALTER TABLE strategy_profile_settings 
ADD [order] INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE strategy_profile_settings 
DROP COLUMN [order];
-- +goose StatementEnd
