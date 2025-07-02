-- +goose Up
-- +goose StatementBegin
INSERT INTO timezones (name, offset, display_name) VALUES
    ('Asia/Riyadh', '+03:00', 'Riyadh');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM timezones WHERE name = 'Asia/Riyadh';
-- +goose StatementEnd
