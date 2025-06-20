-- +goose Up
-- +goose StatementBegin

DECLARE @type_id INT;

-- Identity Proof
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Identity Proof';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'Passport', 'جواز السفر', 'active'),
(@type_id, 'National ID', 'الهوية الوطنية', 'active'),
(@type_id, 'Emirates ID', 'هوية الإمارات', 'active'),
(@type_id, 'Iqama', 'الإقامة', 'active'),
(@type_id, 'Civil ID', 'الهوية المدنية', 'active');

-- Residency Proof
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Residency Proof';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'Residence Permit', 'تصريح الإقامة', 'active'),
(@type_id, 'Emirates ID', 'هوية الإمارات', 'active');

-- Address Proof
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Address Proof';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'Utility Bill', 'فاتورة خدمات', 'active'),
(@type_id, 'Tenancy Contract', 'عقد الإيجار', 'active');

-- Work Authorization
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Work Authorization';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'Work Visa', 'تأشيرة العمل', 'active'),
(@type_id, 'Work Permit', 'تصريح العمل', 'active'),
(@type_id, 'Labour Card', 'بطاقة العمل', 'active');

-- Educational Certificates
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Educational Certificates';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'High School Certificate', 'شهادة الثانوية', 'active'),
(@type_id, 'Degree Certificate', 'شهادة جامعية', 'active'),
(@type_id, 'Equivalency Certificate', 'شهادة المعادلة', 'active');

-- Experience Proof
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Experience Proof';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'Experience Certificate', 'شهادة الخبرة', 'active'),
(@type_id, 'Reference Letter', 'خطاب التوصية', 'active'),
(@type_id, 'Relieving Letter', 'خطاب الإخلاء', 'active');

-- Employment Documents
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Employment Documents';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'Offer Letter', 'عرض العمل', 'active'),
(@type_id, 'Contract Agreement', 'عقد العمل', 'active'),
(@type_id, 'Salary Certificate', 'شهادة الراتب', 'active'),
(@type_id, 'Job Description', 'الوصف الوظيفي', 'active');

-- Legal Compliance
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Legal Compliance';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'NDA', 'اتفاقية عدم الإفشاء', 'active'),
(@type_id, 'Background Check Clearance', 'تصريح الفحص الأمني', 'active'),
(@type_id, 'Security Clearance', 'تصريح الأمان', 'active');

-- Medical Records
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Medical Records';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'Medical Fitness Report', 'تقرير اللياقة الطبية', 'active'),
(@type_id, 'Vaccination Certificate', 'شهادة التطعيم', 'active'),
(@type_id, 'Health Insurance Document', 'وثيقة التأمين الصحي', 'active');

-- Photographs
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Photographs';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'Passport Size Photo', 'صورة بحجم جواز السفر', 'active'),
(@type_id, 'Visa Photo Format', 'صيغة صورة التأشيرة', 'active');

-- Training & Skills
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Training & Skills';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'HSE Certificate', 'شهادة السلامة والصحة المهنية', 'active'),
(@type_id, 'Safety Training', 'تدريب السلامة', 'active'),
(@type_id, 'ISO Certification', 'شهادة الأيزو', 'active'),
(@type_id, 'Language Proficiency', 'إثبات إتقان اللغة', 'active');

-- Travel Documents
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Travel Documents';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'Exit/Re-entry Permit', 'تصريح الخروج والعودة', 'active'),
(@type_id, 'Multiple Entry Visa', 'تأشيرة دخول متعددة', 'active'),
(@type_id, 'Travel History Report', 'تقرير سجل السفر', 'active');

-- Government Registrations
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Government Registrations';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'GOSI Certificate', 'شهادة التأمينات الاجتماعية', 'active'),
(@type_id, 'MOHRE ID', 'بطاقة وزارة العمل', 'active'),
(@type_id, 'Tax ID', 'الرقم الضريبي', 'active'),
(@type_id, 'Social Security Registration', 'تسجيل الضمان الاجتماعي', 'active');

-- Others
SELECT @type_id = id FROM user_document_types WHERE title_en = 'Others';
INSERT INTO user_document_sub_types (user_document_type_id, title_en, title_ar, status) VALUES
(@type_id, 'Miscellaneous', 'متنوع', 'active'),
(@type_id, 'Custom Upload', 'تحميل مخصص', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE user_document_sub_types;
-- +goose StatementEnd
