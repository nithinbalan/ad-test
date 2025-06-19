-- +goose Up
-- +goose StatementBegin
INSERT INTO designations (title_en, title_ar, status) VALUES
('Admin', 'مسؤول', 'active'),
('Chief Executive Officer (CEO)', 'الرئيس التنفيذي', 'active'),
('Chief Financial Officer (CFO)', 'المدير المالي', 'active'),
('Chief Operating Officer (COO)', 'المدير التنفيذي للعمليات', 'active'),
('Chief Technology Officer (CTO)', 'المدير التنفيذي للتكنولوجيا', 'active'),
('Chief Information Officer (CIO)', 'المدير التنفيذي للمعلومات', 'active'),
('Chief Compliance Officer (CCO)', 'المدير التنفيذي للامتثال', 'active'),
('Chief Risk Officer (CRO)', 'المدير التنفيذي للمخاطر', 'active'),
('Chief Audit Executive (CAE)', 'المدير التنفيذي للتدقيق', 'active'),
('Chief Information Security Officer (CISO)', 'المدير التنفيذي لأمن المعلومات', 'active'),
('Vice President of Operations', 'نائب رئيس العمليات', 'active'),
('Vice President of Finance', 'نائب رئيس المالية', 'active'),
('Vice President of Human Resources', 'نائب رئيس الموارد البشرية', 'active'),
('Vice President of IT', 'نائب رئيس تقنية المعلومات', 'active'),
('Director of Compliance', 'مدير الامتثال', 'active'),
('Director of Risk Management', 'مدير إدارة المخاطر', 'active'),
('Director of Internal Audit', 'مدير التدقيق الداخلي', 'active'),
('Director of Information Security', 'مدير أمن المعلومات', 'active'),
('Compliance Manager', 'مدير الامتثال', 'active'),
('Risk Manager', 'مدير المخاطر', 'active'),
('Internal Audit Manager', 'مدير التدقيق الداخلي', 'active'),
('Information Security Manager', 'مدير أمن المعلومات', 'active'),
('Compliance Analyst', 'محلل الامتثال', 'active'),
('Risk Analyst', 'محلل المخاطر', 'active'),
('Internal Auditor', 'مدقق داخلي', 'active'),
('Information Security Analyst', 'محلل أمن المعلومات', 'active'),
('Legal Counsel', 'مستشار قانوني', 'active'),
('Head of Legal', 'رئيس الشؤون القانونية', 'active'),
('Head of Compliance', 'رئيس الامتثال', 'active'),
('Head of Risk Management', 'رئيس إدارة المخاطر', 'active'),
('Head of Internal Audit', 'رئيس التدقيق الداخلي', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE designations;
-- +goose StatementEnd
