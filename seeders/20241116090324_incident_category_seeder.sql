-- +goose Up
-- +goose StatementBegin
INSERT INTO incident_categories (title_en, title_ar, status) VALUES
('Operational Incident', 'حادث تشغيلي', 'active'),
('Security Breach', 'خرق أمني', 'active'),
('Compliance Violation', 'انتهاك الامتثال', 'active'),
('Data Loss', 'فقدان البيانات', 'active'),
('Health and Safety', 'الصحة والسلامة', 'active'),
('Environmental Incident', 'حادث بيئي', 'active'),
('Equipment Failure', 'فشل المعدات', 'active'),
('Network Outage', 'انقطاع الشبكة', 'active'),
('Fraud', 'احتيال', 'active'),
('Customer Complaint', 'شكوى العميل', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE incident_categories;
-- +goose StatementEnd
