-- +goose Up
-- +goose StatementBegin
INSERT INTO email_template_settings (header_image_url, support_text, closing_text, footer_logo_url, footer_address, footer_email, created_by_id, created_at)
VALUES (
    'https://grc-beta.ervarioo.com/assets/email-template/mail-header.png', 
    'For any assistance or questions, feel free to reach out to our support team at info@jethur.com or +61 7536563635.', 
    'Best regards, Jethur GRC Team.', 
    'https://grc-beta.ervarioo.com/assets/email-template/logo.png', 
    'Dammam City, P.O.Box: 76296, Mutlaq Dist.', 
    'Email us at info@jethur.com.', 
    1, 
    GETDATE())
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM email_template_settings WHERE created_by_id = 1;
-- +goose StatementEnd
