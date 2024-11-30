-- +goose Up
-- +goose StatementBegin
INSERT INTO incident_types (title_en, title_ar, status) VALUES
('Information Security Incident', 'حادث أمن المعلومات', 'active'),
('Compliance Breach', 'خرق الامتثال', 'active'),
('Operational Disruption', 'اضطراب تشغيلي', 'active'),
('Data Privacy Incident', 'حادث خصوصية البيانات', 'active'),
('Health and Safety Incident', 'حادث الصحة والسلامة', 'active'),
('Environmental Incident', 'حادث بيئي', 'active'),
('Financial Irregularity', 'مخالفة مالية', 'active'),
('Cybersecurity Attack', 'هجوم إلكتروني', 'active'),
('Physical Security Incident', 'حادث أمني مادي', 'active'),
('Customer Service Issue', 'مشكلة خدمة العملاء', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE incident_types;
-- +goose StatementEnd
