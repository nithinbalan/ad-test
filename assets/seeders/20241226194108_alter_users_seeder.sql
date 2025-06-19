-- +goose Up
-- +goose StatementBegin
DECLARE @UserID INT;
DECLARE @DefaultTimezoneID INT;

-- Retrieve the user ID for the Super Admin
SELECT @UserID = id FROM users WHERE email = 'admin@jethur.com';

-- Retrieve the default timezone ID
SELECT @DefaultTimezoneID = id FROM timezones WHERE name = 'Asia/Riyadh';

-- Insert default user settings for the Super Admin
INSERT INTO user_settings (
    user_id,
    timezone_id,
    language,
    timeout,
    is_notification,
    updated_at,
    updated_by_id
) VALUES (
    @UserID,
    @DefaultTimezoneID,
    'en', -- Default language
    0, -- Default timeout in minutes
    1, -- Enable notifications
    GETDATE(),
    1 -- Admin user ID as updater
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @UserID INT;

-- Retrieve the user ID for the Super Admin
SELECT @UserID = id FROM users WHERE email = 'admin@jethur.com';

-- Delete the user settings for the Super Admin
DELETE FROM user_settings WHERE user_id = @UserID;
-- +goose StatementEnd
