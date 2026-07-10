-- +goose Up
-- +goose StatementBegin
UPDATE modules
SET is_enable = 0
WHERE identifier = 'bcm_report';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules
SET is_enable = 1
WHERE identifier = 'bcm_report';
-- +goose StatementEnd
