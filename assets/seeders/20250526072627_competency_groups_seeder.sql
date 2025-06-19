-- +goose Up
-- +goose StatementBegin
INSERT INTO competency_groups (title_en, title_ar, status) VALUES
('Technical Skills', 'المهارات التقنية', 'active'),
('Frontend Development', 'تطوير الواجهة الأمامية', 'active'),
('Backend Development', 'تطوير الخلفية', 'active'),
('Database Management', 'إدارة قواعد البيانات', 'active'),
('DevOps & CI/CD', 'ديفوبس والتكامل المستمر', 'active'),
('Cloud Infrastructure', 'البنية التحتية السحابية', 'active'),
('API Integration', 'تكامل واجهات برمجة التطبيقات', 'active'),
('Testing & QA', 'الاختبار وضمان الجودة', 'active'),
('AI & Automation', 'الذكاء الاصطناعي والأتمتة', 'active'),
('Mobile Development', 'تطوير التطبيقات المحمولة', 'active'),
('Cybersecurity & Compliance', 'الأمن السيبراني والامتثال', 'active'),
('Threat Management', 'إدارة التهديدات', 'active'),
('SIEM & Security Tools', 'أدوات SIEM والأمان', 'active'),
('Access Control & IAM', 'التحكم في الوصول وإدارة الهوية', 'active'),
('Data Protection', 'حماية البيانات', 'active'),
('Compliance Knowledge', 'معرفة الامتثال', 'active'),
('Secure Development', 'التطوير الآمن', 'active'),
('Risk Management', 'إدارة المخاطر', 'active'),
('Control Management', 'إدارة الضوابط', 'active'),
('Compliance Registers', 'سجلات الامتثال', 'active'),
('Audit Management', 'إدارة التدقيق', 'active'),
('Business Continuity', 'استمرارية الأعمال', 'active'),
('Third-Party Risk', 'مخاطر الأطراف الخارجية', 'active'),
('KPI/KRI Management', 'إدارة مؤشرات الأداء والمخاطر', 'active'),
('Project Planning', 'تخطيط المشاريع', 'active'),
('Agile Methodologies', 'المنهجيات الرشيقة', 'active'),
('Release Management', 'إدارة الإصدارات', 'active'),
('Workflow Automation', 'أتمتة سير العمل', 'active'),
('RACI Matrix Application', 'تطبيق مصفوفة RACI', 'active'),
('Organizational & Leadership', 'القيادة والتنظيم', 'active'),
('Communication & Documentation', 'التواصل والتوثيق', 'active'),
('Sector-Specific', 'حسب القطاع', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE competency_groups;
-- +goose StatementEnd
