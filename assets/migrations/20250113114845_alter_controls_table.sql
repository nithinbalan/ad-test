-- +goose Up
-- +goose StatementBegin
ALTER TABLE controls 
ADD implementation_guideline NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE controls 
DROP COLUMN implementation_guideline;
-- +goose StatementEnd
