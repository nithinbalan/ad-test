-- +goose Up
-- +goose StatementBegin
INSERT INTO timezones (name, offset, display_name) VALUES
    ('Africa/Casablanca', '+00:00', 'Casablanca'),
    ('Africa/Lagos', '+01:00', 'West Central Africa (Lagos)'),
    ('Europe/Istanbul', '+03:00', 'Istanbul'),
    ('Africa/Cairo', '+02:00', 'Cairo'),
    ('Asia/Amman', '+02:00', 'Amman'),
    ('Asia/Jerusalem', '+02:00', 'Jerusalem'),
    ('Asia/Baghdad', '+03:00', 'Baghdad'),
    ('Asia/Kuwait', '+03:00', 'Kuwait'),
    ('Asia/Riyadh', '+03:00', 'Riyadh'),
    ('Europe/Minsk', '+03:00', 'Minsk'),
    ('Asia/Baku', '+04:00', 'Baku'),
    ('Asia/Tbilisi', '+04:00', 'Tbilisi'),
    ('Asia/Yerevan', '+04:00', 'Yerevan'),
    ('Asia/Tashkent', '+05:00', 'Tashkent'),
    ('Asia/Colombo', '+05:30', 'Sri Jayawardenepura'),
    ('Asia/Almaty', '+06:00', 'Almaty'),
    ('Asia/Novosibirsk', '+07:00', 'Novosibirsk'),
    ('Asia/Ulaanbaatar', '+08:00', 'Ulaanbaatar'),
    ('Australia/Darwin', '+09:30', 'Darwin'),
    ('Pacific/Noumea', '+11:00', 'Noumea'),
    ('Pacific/Tarawa', '+12:00', 'Tarawa'),
    ('Pacific/Fakaofo', '+13:00', 'Tokelau'),
    ('Pacific/Kiritimati', '+14:00', 'Kiritimati (Christmas Island)');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM timezones WHERE name IN (
    'Africa/Casablanca',
    'Africa/Lagos',
    'Africa/Cairo',
    'Europe/Istanbul',
    'Asia/Amman',
    'Asia/Jerusalem',
    'Asia/Baghdad',
    'Asia/Kuwait',
    'Asia/Riyadh',
    'Europe/Minsk',
    'Asia/Baku',
    'Asia/Tbilisi',
    'Asia/Yerevan',
    'Asia/Tashkent',
    'Asia/Colombo',
    'Asia/Almaty',
    'Asia/Novosibirsk',
    'Asia/Ulaanbaatar',
    'Australia/Darwin',
    'Pacific/Noumea',
    'Pacific/Tarawa',
    'Pacific/Fakaofo',
    'Pacific/Kiritimati'
);
-- +goose StatementEnd
