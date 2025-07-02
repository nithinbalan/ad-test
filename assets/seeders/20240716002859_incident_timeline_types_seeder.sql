-- +goose Up
-- +goose StatementBegin
INSERT INTO incident_timeline_types
    (title_en, title_ar, color_code, type)
VALUES
    ('Reporting', 'تقرير', '#28A745', 'reporting'),
    ('Investigating', 'تحقيق', '#007BFF', 'investigating'),
    ('RCA', 'تحليل السبب الجذري', '#FFC107', 'rca'),
    ('Communication', 'الاتصالات', '#FD7E14', 'communication'),
    ('Action Plans', 'خطط العمل', '#6F42C1', 'action-plans'),
    ('Closure', 'إغلاق', '#E0554D', 'closure'),
    ('PIR', 'مراجعة ما بعد الحادث', '#DC3545', 'pir');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE incident_timeline_types;
-- +goose StatementEnd

