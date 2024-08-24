-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_timeline_types
    (title_en, title_ar, color_code, type)
VALUES
    ('Initiated', 'تم البدء', '#17A2B8', 'initiated'),
    ('Meetings', 'الاجتماعات', '#007BFF', 'meetings'),
    ('Workpapers', 'أوراق العمل', '#28A745', 'workpapers'),
    ('Request Documents', 'طلب المستندات', '#FFC107', 'request-documents'),
    ('Schedule', 'الجدول', '#6F42C1', 'schedule'),
    ('Checklist', 'قائمة الفحص', '#FD7E14', 'checklist'),
    ('Draft Report', 'تقرير المسودة', '#FFA500', 'draft-report'),
    ('Final Report', 'التقرير النهائي', '#DC3545', 'final-report'),
    ('Publish', 'النشر', '#343A40', 'publish');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_timeline_types;
-- +goose StatementEnd
