-- +goose Up
-- +goose StatementBegin
INSERT INTO email_template_settings (header_image_url, support_text, closing_text, footer_logo_url, footer_address, footer_email, created_by_id, created_at)
VALUES (
    'https://grc-beta.ervarioo.com/assets/email-template/mail-header.png', 
    'For any assistance or questions, feel free to reach out to our support team at', 
    'Best regards, Jethur GRC Team.', 
    'https://grc-beta.ervarioo.com/assets/email-template/logo.png', 
    'Riyadh (HQ), Umar Ibn Abdul Aziz Br Rd, Ar Rabwah, Riyadh 12834', 
    'Email us at info@jethur.com.', 
    1, 
    GETDATE())
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM email_template_settings WHERE created_by_id = 1;
-- +goose StatementEnd
