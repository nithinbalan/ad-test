-- +goose Up
-- +goose StatementBegin
INSERT INTO planned_advisory_statuses (title_en, title_ar, description_en, description_ar, color_code, type, status) VALUES
    ('Open', N'مفتوح', 'Received; no action yet.', N'تم الاستلام؛ لم يتم اتخاذ أي إجراء بعد.', '#007BFF', 'open', 'active'),
    ('Acknowledged', N'تم الإقرار', 'Reviewed and logged; no further action needed.', N'تمت المراجعة والتسجيل؛ لا حاجة لمزيد من الإجراءات.', '#28A745', 'acknowledged', 'active'),
    ('Awaits Feedback', N'في انتظار الرد', 'Reply sent to the external party; awaiting confirmation.', N'تم إرسال الرد إلى الطرف الخارجي؛ في انتظار التأكيد.', '#17A2B8', 'awaits_feedback', 'active'),
    ('Feedback Under Review', N'الرد قيد المراجعة', 'External reply received and under internal review.', N'تم استلام الرد الخارجي وهو قيد المراجعة الداخلية.', '#FFC107', 'feedback_under_review', 'active'),
    ('Awaits Feedback Escalated', N'في انتظار الرد المصعد', 'No reply received and escalated internally.', N'لم يتم استلام رد وتم التصعيد داخليًا.', '#FD7E14', 'awaits_feedback_escalated', 'active'),
    ('Resolved', N'تم الحل', 'Handled and formally closed.', N'تمت المعالجة والإغلاق رسميًا.', '#6C757D', 'resolved', 'active'),
    ('Rejected', N'مرفوض', 'Not applicable or rejected after review.', N'غير قابل للتطبيق أو مرفوض بعد المراجعة.', '#DC3545', 'rejected', 'active'),
    ('Reopened', N'تم إعادة الفتح', 'Reopened after closure.', N'تم إعادة الفتح بعد الإغلاق.', '#FFA500', 'reopened', 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM planned_advisory_statuses
WHERE type IN (
    'open',
    'acknowledged',
    'awaits_feedback',
    'feedback_under_review',
    'awaits_feedback_escalated',
    'resolved',
    'rejected',
    'reopened'
);
-- +goose StatementEnd
