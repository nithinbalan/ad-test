-- +goose Up
-- +goose StatementBegin
INSERT INTO incident_timeline_types
    (title_en, title_ar, color_code, type)
VALUES
    ('Reporting', 'تقرير', '#FF0000', 'reporting'),
    ('Investigating', 'تحقيق', '#FFFF00', 'investigating'),
    ('RCA', 'تحليل السبب الجذري', '#0000FF', 'rca'),
    ('Communication', 'الاتصالات', '#800080', 'communication'),
    ('Action Plans', 'خطط العمل', '#FFA500', 'action-plans'),
    ('Closure', 'إغلاق', '#008000', 'closure'),
    ('PIR', 'مراجعة ما بعد الحادث', '#00FFFF', 'pir');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE incident_timeline_types;
-- +goose StatementEnd

