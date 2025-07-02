-- +goose Up
-- +goose StatementBegin
INSERT INTO feedback_options (title_en, title_ar, color_code, type, status) VALUES
('Poor', 'في الوقت المحدد', '#DC3545', 'poor', 'active'),
('Fair', 'تأخير', '#FFC107', 'fair', 'active'),
('Good', 'مبكر', '#17A2B8', 'good', 'active'),
('Excellent', 'مبكر', '#90DD75', 'excellent', 'active');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE feedback_options;
-- +goose StatementEnd