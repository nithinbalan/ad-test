-- +goose Up
-- +goose StatementBegin
INSERT INTO timezones (name, offset, display_name) VALUES
('Pacific Standard', '-08:00', 'Pacific Time (US & Canada)'),
('Eastern Standard', '-05:00', 'Eastern Time (US & Canada)'),
('Greenwich Mean Time', '+00:00', 'Greenwich Mean Time'),
('Central European', '+01:00', 'Central European Time (Berlin)'),
('India Standard', '+05:30', 'India Time (New Delhi, Mumbai)');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE timezones;
-- +goose StatementEnd
