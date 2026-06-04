-- +goose Up
-- +goose StatementBegin
-- Delete all business impact assessments before updating impact categories

-- Update old existing categories with new type values
UPDATE impact_categories
SET type = 'financial'
WHERE title_en = 'Financial';

UPDATE impact_categories
SET type = 'reputational'
WHERE title_en = 'Reputational';

UPDATE impact_categories
SET type = 'operational'
WHERE title_en = 'Operational';

UPDATE impact_categories
SET type = 'human-resource'
WHERE title_en = 'Human Resources';

-- Insert new categories
INSERT INTO impact_categories (title_en, title_ar, type, status)
SELECT 'Downtime', 'توقف', 'downtime', 'active'
WHERE NOT EXISTS (SELECT 1 FROM impact_categories WHERE title_en = 'Downtime');

INSERT INTO impact_categories (title_en, title_ar, type, status)
SELECT 'Information Technology', 'تكنولوجيا المعلومات', 'information-technology', 'active'
WHERE NOT EXISTS (SELECT 1 FROM impact_categories WHERE title_en = 'Information Technology');

INSERT INTO impact_categories (title_en, title_ar, type, status)
SELECT 'Cyber Security', 'الأمن السيبراني', 'cyber-security', 'active'
WHERE NOT EXISTS (SELECT 1 FROM impact_categories WHERE title_en = 'Cyber Security');

INSERT INTO impact_categories (title_en, title_ar, type, status)
SELECT 'Customer Satisfaction', 'رضا العملاء', 'customer-satisfaction', 'active'
WHERE NOT EXISTS (SELECT 1 FROM impact_categories WHERE title_en = 'Customer Satisfaction');

INSERT INTO impact_categories (title_en, title_ar, type, status)
SELECT 'Development and Growth', 'التطوير والنمو', 'development-and-growth', 'active'
WHERE NOT EXISTS (SELECT 1 FROM impact_categories WHERE title_en = 'Development and Growth');

INSERT INTO impact_categories (title_en, title_ar, type, status)
SELECT 'Insurance', 'التأمين', 'insurance', 'active'
WHERE NOT EXISTS (SELECT 1 FROM impact_categories WHERE title_en = 'Insurance');

INSERT INTO impact_categories (title_en, title_ar, type, status)
SELECT 'Regulations and Policies', 'اللوائح والسياسات', 'regulations-and-policies', 'active'
WHERE NOT EXISTS (SELECT 1 FROM impact_categories WHERE title_en = 'Regulations and Policies');

INSERT INTO impact_categories (title_en, title_ar, type, status)
SELECT 'Health and Safety', 'الصحة والسلامة', 'health-and-safety', 'active'
WHERE NOT EXISTS (SELECT 1 FROM impact_categories WHERE title_en = 'Health and Safety');

INSERT INTO impact_categories (title_en, title_ar, type, status)
SELECT 'Governance and Internal Control', 'الحوكمة والرقابة الداخلية', 'governance-and-internal-control', 'active'
WHERE NOT EXISTS (SELECT 1 FROM impact_categories WHERE title_en = 'Governance and Internal Control');

-- Delete unwanted old categories (Legal/Compliance, IT/System, Supply Chain)
DELETE FROM impact_categories
WHERE title_en IN ('Legal/Compliance', 'IT/System', 'Supply Chain');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert the type back to NULL
UPDATE impact_categories
SET type = NULL
WHERE type = 'financial';

UPDATE impact_categories
SET type = NULL
WHERE type = 'reputational';

UPDATE impact_categories
SET type = NULL
WHERE type = 'operational';

UPDATE impact_categories
SET type = NULL
WHERE type = 'downtime';

UPDATE impact_categories
SET type = NULL
WHERE type = 'information-technology';

UPDATE impact_categories
SET type = NULL
WHERE type = 'cyber-security';

UPDATE impact_categories
SET type = NULL
WHERE type = 'customer-satisfaction';

UPDATE impact_categories
SET type = NULL
WHERE type = 'development-and-growth';

UPDATE impact_categories
SET type = NULL
WHERE type = 'insurance';

UPDATE impact_categories
SET type = NULL
WHERE type = 'human-resource';

UPDATE impact_categories
SET type = NULL
WHERE type = 'regulations-and-policies';

UPDATE impact_categories
SET type = NULL
WHERE type = 'health-and-safety';

UPDATE impact_categories
SET type = NULL
WHERE type = 'governance-and-internal-control';
-- +goose StatementEnd
