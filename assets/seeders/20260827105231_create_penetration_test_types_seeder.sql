-- +goose Up
-- +goose StatementBegin
INSERT INTO penetration_test_types (code, title_en, title_ar, status)
VALUES
    ('PT-01', 'Vulnerability assessment', N'تقييم الثغرات', 'active'),
    ('PT-08', 'Database testing', N'اختبار قواعد البيانات', 'active'),
    ('PT-02', 'Internal web-application pentest', N'اختبار اختراق تطبيقات الويب الداخلية', 'active'),
    ('PT-09', 'Source-code review', N'مراجعة الكود المصدري', 'active'),
    ('PT-03', 'Internal-network pentest', N'اختبار اختراق الشبكة الداخلية', 'active'),
    ('PT-10', 'Configuration review', N'مراجعة الإعدادات', 'active'),
    ('PT-04', 'Wireless-network pentest', N'اختبار اختراق الشبكات اللاسلكية', 'active'),
    ('PT-11', 'Login systems (AD/IAM/SSO)', N'أنظمة تسجيل الدخول (AD/IAM/SSO)', 'active'),
    ('PT-05', 'External-network pentest', N'اختبار اختراق الشبكة الخارجية', 'active'),
    ('PT-12', 'Cybersecurity-systems pentest', N'اختبار اختراق أنظمة الأمن السيبراني', 'active'),
    ('PT-06', 'API security testing', N'اختبار أمان واجهات برمجة التطبيقات', 'active'),
    ('PT-13', 'Social engineering', N'الهندسة الاجتماعية', 'active'),
    ('PT-07', 'Mobile-application testing', N'اختبار تطبيقات الهاتف المحمول', 'active'),
    ('PT-14', 'Physical security testing', N'اختبار الأمن المادي', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE penetration_test_types;
-- +goose StatementEnd
