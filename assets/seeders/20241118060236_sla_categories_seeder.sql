-- +goose Up
-- +goose StatementBegin

INSERT INTO sla_categories (title_en, title_ar, status)
VALUES
    ('Response Time', 'زمن الاستجابة', 'active'),
    ('Resolution Time', 'وقت الحل', 'active'),
    ('Uptime Guarantee', 'ضمان وقت التشغيل', 'active'),
    ('Service Availability', 'توفر الخدمة', 'active'),
    ('Incident Management', 'إدارة الحوادث', 'active'),
    ('Problem Management', 'إدارة المشاكل', 'active'),
    ('Change Management', 'إدارة التغيير', 'active'),
    ('Data Security', 'أمان البيانات', 'active'),
    ('Compliance Monitoring', 'مراقبة الامتثال', 'active'),
    ('Reporting and Analytics', 'التقارير والتحليلات', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

TRUNCATE TABLE sla_categories;

-- +goose StatementEnd
