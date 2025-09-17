-- +goose Up
-- +goose StatementBegin
ALTER TABLE meetings ADD transcript NVARCHAR(MAX);
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
ALTER TABLE meetings DROP COLUMN transcript;
-- +goose StatementEnd