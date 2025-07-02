-- +goose Up
-- +goose StatementBegin
INSERT INTO test_result_options (title_en, title_ar, color_code, type, status) VALUES
('Successful', 'في الوقت المحدد', '#90DD75', 'successful', 'active'),
('Partially Successful', 'تأخير', '#17A2B8', 'partially-successful', 'active'),
('Failure', 'مبكر', '#DC3545', 'failure', 'active');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE test_result_options;
-- +goose StatementEnd