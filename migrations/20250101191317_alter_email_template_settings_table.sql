-- +goose Up
-- +goose StatementBegin
ALTER TABLE email_template_settings
ALTER COLUMN header_image_url NVARCHAR(MAX) NULL;

ALTER TABLE email_template_settings
ALTER COLUMN footer_logo_url NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE email_template_settings
ALTER COLUMN header_image_url NVARCHAR(MAX) NOT NULL;

ALTER TABLE email_template_settings
ALTER COLUMN footer_logo_url NVARCHAR(MAX) NOT NULL;
-- +goose StatementEnd
