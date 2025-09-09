-- +goose Up
-- +goose StatementBegin
UPDATE modules
SET [order] = 17
WHERE identifier = 'policy_management';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules 
SET [order] = 16 
WHERE identifier = 'policy_management';
-- +goose StatementEnd
