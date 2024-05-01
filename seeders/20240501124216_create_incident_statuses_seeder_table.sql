-- +goose Up
-- +goose StatementBegin
INSERT INTO incident_statuses
    (title_en, title_ar, description_en, description_ar, color_code, type, status)
VALUES
    ('Open', 'مفتوح', 'The incident is open and under investigation.', 'الحادث مفتوح وقيد التحقيق.', '#FF0000', 'open', 'active'),
    ('Investigating', 'جارٍ التحقيق', 'The incident is under investigation to determine its cause and severity.', 'الحادث قيد التحقيق لتحديد سببه وشدته.', '#FFFF00', 'investigating', 'active'),
    ('In Progress', 'في تَقَدم', 'The incident is in progress and actions are being taken to resolve it.', 'الحادث قيد التقدم ويتم اتخاذ إجراءات لحله.', '#FFA500', 'in_progress', 'active'),
    ('Resolved', 'محلول', 'The incident has been resolved and no longer requires attention.', 'تم حل الحادث ولا يتطلب المزيد من الاهتمام.', '#00FF00', 'resolved', 'active'),
    ('Closed', 'مغلق', 'The incident has been closed after resolution and review.', 'تم إغلاق الحادث بعد حله ومراجعته.', '#808080', 'closed', 'active'),
    ('Pending Review', 'في انتظار المراجعة', 'The incident is pending review by appropriate authorities.', 'الحادث في انتظار المراجعة من قبل الجهات المختصة.', '#800080', 'pending_review', 'active'),
    ('Escalated', 'تم تصعيده', 'The incident has been escalated to higher management for attention.', 'تم تصعيد الحادث إلى إدارة أعلى للتنبه إليه.', '#0000FF', 'escalated', 'active'),
    ('Re-opened', 'إعادة فتح', 'The incident has been reopened for further investigation or action.', 'تم إعادة فتح الحادث لإجراء تحقيق أو إجراء إضافي.', '#FF00FF', 're_opened', 'active'),
    ('Deferred', 'مؤجل', 'The incident has been deferred to be addressed later.', 'تم تأجيل الحادث للتعامل معه لاحقاً.', '#00FFFF', 'deferred', 'active'),
    ('Mitigated', 'تم تخفيفه', 'The incident has been mitigated to reduce its impact.', 'تم تخفيف الحادث لتقليل تأثيره.', '#008000', 'mitigated', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE incident_statuses;
-- +goose StatementEnd
