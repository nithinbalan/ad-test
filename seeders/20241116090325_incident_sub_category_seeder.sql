-- +goose Up
-- +goose StatementBegin

-- Declare variables for specific incident category IDs
DECLARE @operationalIncidentID INT;
DECLARE @securityBreachID INT;

SELECT @operationalIncidentID = id FROM incident_categories WHERE title_en = 'Operational Incident';

SELECT @securityBreachID = id FROM incident_categories WHERE title_en = 'Security Breach';

INSERT INTO incident_sub_categories (incident_category_id, title_en, title_ar, status)
VALUES
    (@operationalIncidentID, 'System Downtime', 'توقف النظام', 'active'),
    (@operationalIncidentID, 'Process Disruption', 'اضطراب العملية', 'active'),
    (@operationalIncidentID, 'Supply Chain Delay', 'تأخير سلسلة التوريد', 'active'),
    (@operationalIncidentID, 'Power Outage', 'انقطاع الكهرباء', 'active'),
    (@operationalIncidentID, 'Resource Unavailability', 'عدم توفر الموارد', 'active');

INSERT INTO incident_sub_categories (incident_category_id, title_en, title_ar, status)
VALUES
    (@securityBreachID, 'Unauthorized Access', 'دخول غير مصرح به', 'active'),
    (@securityBreachID, 'Data Breach', 'اختراق البيانات', 'active'),
    (@securityBreachID, 'Malware Attack', 'هجوم البرمجيات الخبيثة', 'active'),
    (@securityBreachID, 'Phishing Attack', 'هجوم التصيد', 'active'),
    (@securityBreachID, 'Denial of Service', 'هجوم حجب الخدمة', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE incident_sub_categories;

-- +goose StatementEnd
