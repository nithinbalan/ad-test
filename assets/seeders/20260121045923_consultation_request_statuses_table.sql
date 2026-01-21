
-- +goose Up
-- +goose StatementBegin
INSERT INTO consultation_request_statuses (title_en, title_ar, description_en, description_ar, type, status, created_at) VALUES
('New', 'جديد', 'New consultation request', 'طلب استشارة جديد', 'new', 'active', GETDATE()),
('Needs Clarification', 'يحتاج توضيح', 'Request needs clarification', 'الطلب يحتاج توضيح', 'needs_clarification', 'active', GETDATE()),
('Under Processing', 'قيد المعالجة', 'Request is under processing', 'الطلب قيد المعالجة', 'under_processing', 'active', GETDATE()),
('Reply Sent', 'تم إرسال الرد', 'Reply has been sent', 'تم إرسال الرد', 'reply_sent', 'active', GETDATE()),
('Closed', 'مغلق', 'Request is closed', 'الطلب مغلق', 'closed', 'active', GETDATE()),
('Reopened', 'أعيد فتحه', 'Request has been reopened', 'تم إعادة فتح الطلب', 'reopened', 'active', GETDATE());
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE consultation_request_statuses;
-- +goose StatementEnd
