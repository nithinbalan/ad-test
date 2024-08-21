-- +goose Up
-- +goose StatementBegin

-- Add new columns
ALTER TABLE user_notifications
ADD
    notification_type NVARCHAR(255) NOT NULL,
    reference_number VARCHAR(45) NOT NULL,
    heading_en NVARCHAR(255) NOT NULL,
    heading_ar NVARCHAR(255) NOT NULL,
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    url_param NVARCHAR(MAX) NOT NULL,
    query_param NVARCHAR(MAX) NOT NULL,
    is_system_notification BIT DEFAULT 0;

-- Make the created_by_id column nullable
ALTER TABLE user_notifications
ALTER COLUMN created_by_id INT NULL;

-- Drop old columns that are no longer needed
ALTER TABLE user_notifications
DROP COLUMN 
    icon,
    message_en,
    message_ar,
    redirect_url;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Add the old columns back
ALTER TABLE user_notifications
ADD
    icon NVARCHAR(45) NOT NULL,
    message_en NVARCHAR(255) NOT NULL,
    message_ar NVARCHAR(255) NOT NULL,
    redirect_url NVARCHAR(255) NOT NULL;

-- Revert the created_by_id column to NOT NULL
ALTER TABLE user_notifications
ALTER COLUMN created_by_id INT NOT NULL;

-- Drop the new columns that were added
ALTER TABLE user_notifications
DROP COLUMN 
    notification_type,
    reference_number,
    heading_en,
    heading_ar,
    title_en,
    title_ar,
    url_param,
    query_param,
    is_system_notification;

-- +goose StatementEnd
