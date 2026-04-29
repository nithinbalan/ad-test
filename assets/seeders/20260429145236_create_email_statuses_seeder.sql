-- +goose Up
-- +goose StatementBegin
INSERT INTO email_statuses 
(title_en, title_ar, color_code, [type]) 
VALUES
('Pending', 'قيد الانتظار', '#FFC107', 'pending'),
('Processing', 'قيد المعالجة', '#17A2B8', 'processing'),
('Sent', 'تم الإرسال', '#28A745', 'sent'),
('Failed', 'فشل الإرسال', '#DC3545', 'failed');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM email_statuses;
-- +goose StatementEnd