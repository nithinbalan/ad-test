-- +goose Up
-- +goose StatementBegin
 
INSERT INTO user_document_types (title_en, title_ar, status) VALUES
('Identity Proof', 'إثبات الهوية', 'active'),
('Residency Proof', 'إثبات الإقامة', 'active'),
('Address Proof', 'إثبات العنوان', 'active'),
('Work Authorization', 'تفويض العمل', 'active'),
('Educational Certificates', 'شهادات التعليم', 'active'),
('Experience Proof', 'إثبات الخبرة', 'active'),
('Employment Documents', 'مستندات التوظيف', 'active'),
('Legal Compliance', 'الامتثال القانوني', 'active'),
('Medical Records', 'السجلات الطبية', 'active'),
('Photographs', 'الصور', 'active'),
('Training & Skills', 'التدريب والمهارات', 'active'),
('Travel Documents', 'مستندات السفر', 'active'),
('Government Registrations', 'التسجيلات الحكومية', 'active'),
('Others', 'أخرى', 'active');
 
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE user_document_types;
-- +goose StatementEnd