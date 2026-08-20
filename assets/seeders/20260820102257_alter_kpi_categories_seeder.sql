-- +goose Up
-- +goose StatementBegin
UPDATE kpi_categories SET title_en = 'Acceptable Use Policy', title_ar = N'سياسة الاستخدام المقبول', type = 'acceptable-use-policy' WHERE title_en = 'Financial';
UPDATE kpi_categories SET title_en = 'Application Whitelisting', title_ar = N'القائمة البيضاء للتطبيقات', type = 'application-whitelisting' WHERE title_en = 'Operational';
UPDATE kpi_categories SET title_en = 'Asset Classification', title_ar = N'تصنيف الأصول', type = 'asset-classification' WHERE title_en = 'Compliance';
UPDATE kpi_categories SET title_en = 'Asset Discovery', title_ar = N'اكتشاف الأصول', type = 'asset-discovery' WHERE title_en = 'Customer';

INSERT INTO kpi_categories (title_en, title_ar, status, type)
VALUES
    ('Asset Maintenance', N'صيانة الأصول', 'active', 'asset-maintenance'),
    ('Backup & Recovery Management', N'إدارة النسخ الاحتياطي والاسترداد', 'active', 'backup-recovery-management'),
    ('Bring Your Own Device (BYOD)', N'أحضر جهازك الخاص', 'active', 'byod'),
    ('Change Management', N'إدارة التغيير', 'active', 'change-management'),
    ('Cloud Services', N'خدمات سحابية', 'active', 'cloud-services'),
    ('Configuration Management & Hardening', N'إدارة التكوين والتقوية', 'active', 'configuration-management-hardening'),
    ('Cryptography', N'التشفير', 'active', 'cryptography'),
    ('Customer Cybersecurity Awareness', N'الوعي بالأمن السيبراني للعملاء', 'active', 'customer-cybersecurity-awareness'),
    ('Cybersecurity Audit', N'تدقيق الأمن السيبراني', 'active', 'cybersecurity-audit'),
    ('Cybersecurity Awareness', N'الوعي بالأمن السيبراني', 'active', 'cybersecurity-awareness'),
    ('Cybersecurity Awareness & Training', N'التوعية والتدريب في مجال الأمن السيبراني', 'active', 'cybersecurity-awareness-training'),
    ('Cybersecurity Compliance', N'الامتثال للأمن السيبراني', 'active', 'cybersecurity-compliance'),
    ('Cybersecurity in Human Resources', N'الأمن السيبراني في الموارد البشرية', 'active', 'cybersecurity-hr'),
    ('Cybersecurity in Project Management', N'الأمن السيبراني في إدارة المشاريع', 'active', 'cybersecurity-project-management'),
    ('Cybersecurity Management', N'إدارة الأمن السيبراني', 'active', 'cybersecurity-management'),
    ('Cybersecurity Risk Assessment', N'تقييم مخاطر الأمن السيبراني', 'active', 'cybersecurity-risk-assessment'),
    ('Cybersecurity Risk Treatment & Monitoring', N'معالجة ومراقبة مخاطر الأمن السيبراني', 'active', 'cybersecurity-risk-treatment-monitoring'),
    ('Cybersecurity Strategy', N'استراتيجية الأمن السيبراني', 'active', 'cybersecurity-strategy'),
    ('Email & Web Browser Protection', N'حماية البريد الإلكتروني ومتصفح الويب', 'active', 'email-web-browser-protection'),
    ('Identity & Access Management', N'إدارة الهوية والوصول', 'active', 'identity-access-management'),
    ('Incident Management', N'إدارة الحوادث', 'active', 'incident-management'),
    ('Information Protection', N'حماية المعلومات', 'active', 'information-protection'),
    ('Logging & Event Management', N'التسجيل وإدارة الأحداث', 'active', 'logging-event-management'),
    ('Malware Handling', N'التعامل مع البرامج الضارة', 'active', 'malware-handling'),
    ('Network Security', N'أمن الشبكات', 'active', 'network-security'),
    ('Outsourcing Services', N'خدمات التعهيد', 'active', 'outsourcing-services'),
    ('Patch Management', N'إدارة التصحيح', 'active', 'patch-management'),
    ('Penetration Testing', N'اختبار الاختراق', 'active', 'penetration-testing'),
    ('Physical Access Management', N'إدارة الوصول المادي', 'active', 'physical-access-management'),
    ('Protection of Physical Information Assets', N'حماية أصول المعلومات المادية', 'active', 'physical-information-protection'),
    ('Secure Disposal of Assets', N'التخلص الآمن من الأصول', 'active', 'secure-disposal-assets'),
    ('Secure Software Development', N'تطوير البرمجيات الآمنة', 'active', 'secure-software-development'),
    ('Vulnerability Management', N'إدارة الثغرات الأمنية', 'active', 'vulnerability-management');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM kpi_categories
WHERE type IN (
    'asset-maintenance', 'backup-recovery-management', 'byod', 'change-management',
    'cloud-services', 'configuration-management-hardening', 'cryptography',
    'customer-cybersecurity-awareness', 'cybersecurity-audit', 'cybersecurity-awareness',
    'cybersecurity-awareness-training', 'cybersecurity-compliance', 'cybersecurity-hr',
    'cybersecurity-project-management', 'cybersecurity-management',
    'cybersecurity-risk-assessment', 'cybersecurity-risk-treatment-monitoring',
    'cybersecurity-strategy', 'email-web-browser-protection', 'identity-access-management',
    'incident-management', 'information-protection', 'logging-event-management',
    'malware-handling', 'network-security', 'outsourcing-services', 'patch-management',
    'penetration-testing', 'physical-access-management', 'physical-information-protection',
    'secure-disposal-assets', 'secure-software-development', 'vulnerability-management'
);

UPDATE kpi_categories SET title_en = 'Financial', title_ar = N'مالي', type = 'other' WHERE type = 'acceptable-use-policy';
UPDATE kpi_categories SET title_en = 'Operational', title_ar = N'التشغيلية', type = 'other' WHERE type = 'application-whitelisting';
UPDATE kpi_categories SET title_en = 'Compliance', title_ar = N'امتثال', type = 'other' WHERE type = 'asset-classification';
UPDATE kpi_categories SET title_en = 'Customer', title_ar = N'عميل', type = 'other' WHERE type = 'asset-discovery';
-- +goose StatementEnd
