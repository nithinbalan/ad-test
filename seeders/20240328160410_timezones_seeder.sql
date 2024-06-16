-- +goose Up
-- +goose StatementBegin
INSERT INTO timezones (name, offset, display_name) VALUES
('Pacific/Midway', '-11:00', 'Midway Island, Samoa'),
('Pacific/Honolulu', '-10:00', 'Hawaii'),
('America/Anchorage', '-09:00', 'Alaska'),
('America/Los_Angeles', '-08:00', 'Pacific Time (US & Canada)'),
('America/Denver', '-07:00', 'Mountain Time (US & Canada)'),
('America/Chicago', '-06:00', 'Central Time (US & Canada)'),
('America/New_York', '-05:00', 'Eastern Time (US & Canada)'),
('America/Caracas', '-04:00', 'Caracas'),
('America/Halifax', '-04:00', 'Atlantic Time (Canada)'),
('America/Sao_Paulo', '-03:00', 'Brasilia'),
('Atlantic/South_Georgia', '-02:00', 'Mid-Atlantic'),
('Atlantic/Azores', '-01:00', 'Azores'),
('Europe/London', '+00:00', 'Greenwich Mean Time'),
('Europe/Berlin', '+01:00', 'Central European Time (Berlin)'),
('Europe/Kaliningrad', '+02:00', 'Kaliningrad'),
('Europe/Moscow', '+03:00', 'Moscow, St. Petersburg, Volgograd'),
('Asia/Tehran', '+03:30', 'Tehran'),
('Asia/Dubai', '+04:00', 'Abu Dhabi, Muscat'),
('Asia/Kabul', '+04:30', 'Kabul'),
('Asia/Karachi', '+05:00', 'Islamabad, Karachi'),
('Asia/Kolkata', '+05:30', 'Chennai, Kolkata, Mumbai, New Delhi'),
('Asia/Kathmandu', '+05:45', 'Kathmandu'),
('Asia/Dhaka', '+06:00', 'Dhaka'),
('Asia/Rangoon', '+06:30', 'Yangon (Rangoon)'),
('Asia/Bangkok', '+07:00', 'Bangkok, Hanoi, Jakarta'),
('Asia/Shanghai', '+08:00', 'Beijing, Chongqing, Hong Kong, Urumqi'),
('Asia/Tokyo', '+09:00', 'Osaka, Sapporo, Tokyo'),
('Australia/Adelaide', '+09:30', 'Adelaide'),
('Australia/Sydney', '+10:00', 'Canberra, Melbourne, Sydney'),
('Pacific/Guadalcanal', '+11:00', 'Solomon Is.'),
('Pacific/Fiji', '+12:00', 'Fiji, Kamchatka, Marshall Is.');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE timezones;
-- +goose StatementEnd
