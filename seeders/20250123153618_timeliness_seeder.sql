-- +goose Up
-- +goose StatementBegin
INSERT INTO timeliness (title_en, title_ar, color_code, type, status) VALUES
('Ontime', 'في الوقت المحدد', '#90DD75', 'ontime', 'active'),
('Delayed', 'تأخير', '#DC3545', 'delayed', 'active'),
('Early', 'مبكر', '#17A2B8', 'early', 'active');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE timeliness;
-- +goose StatementEnd