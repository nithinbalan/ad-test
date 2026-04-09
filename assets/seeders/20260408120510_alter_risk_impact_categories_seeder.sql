-- +goose Up
-- +goose StatementBegin
-- Update existing common categories to match new naming
UPDATE risk_impact_categories SET title_en = 'Financial', title_ar = 'المالي' WHERE title_en = 'Financial Impact';
UPDATE risk_impact_categories SET title_en = 'Operational', title_ar = 'التشغيلي' WHERE title_en = 'Operational Impact';
UPDATE risk_impact_categories SET title_en = 'Reputation', title_ar = 'السمعة' WHERE title_en = 'Reputational Impact';
UPDATE risk_impact_categories SET title_en = 'Health and Safety', title_ar = 'الصحة والسلامة' WHERE title_en = 'Health and Safety Impact';

-- Delete categories that are no longer needed
DELETE FROM risk_impact_categories WHERE title_en IN ('Compliance Impact', 'Legal Impact', 'Environmental Impact', 'Customer Impact', 'Technological Impact', 'Strategic Impact');

-- Insert new categories
INSERT INTO risk_impact_categories (title_en, title_ar, status) VALUES
('Business Interruption', 'انقطاع الأعمال', 'active'),
('Technical', 'التقني', 'active'),
('Cyber', 'السيبراني', 'active'),
('Customer Satisfaction', 'رضا العملاء', 'active'),
('Development and Growth', 'التطوير والنمو', 'active'),
('Credit', 'الائتمان', 'active'),
('Human Resources', 'الموارد البشرية', 'active'),
('Regulations and Policies', 'اللوائح والسياسات', 'active'),
('Governance and Internal Controls', 'الحوكمة والضوابط الداخلية', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert updates to original names
UPDATE risk_impact_categories SET title_en = 'Financial Impact', title_ar = 'التأثير المالي' WHERE title_en = 'Financial';
UPDATE risk_impact_categories SET title_en = 'Operational Impact', title_ar = 'التأثير التشغيلي' WHERE title_en = 'Operational';
UPDATE risk_impact_categories SET title_en = 'Reputational Impact', title_ar = 'التأثير على السمعة' WHERE title_en = 'Reputation';
UPDATE risk_impact_categories SET title_en = 'Health and Safety Impact', title_ar = 'التأثير على الصحة والسلامة' WHERE title_en = 'Health and Safety';

-- Delete new categories
DELETE FROM risk_impact_categories WHERE title_en IN ('Business Interruption', 'Technical', 'Cyber', 'Customer Satisfaction', 'Development and Growth', 'Credit', 'Human Resources', 'Regulations and Policies', 'Governance and Internal Controls');

-- Restore deleted categories
INSERT INTO risk_impact_categories (title_en, title_ar, status) VALUES
('Compliance Impact', 'تأثير الامتثال', 'active'),
('Legal Impact', 'التأثير القانوني', 'active'),
('Environmental Impact', 'التأثير البيئي', 'active'),
('Customer Impact', 'التأثير على العملاء', 'active'),
('Technological Impact', 'التأثير التكنولوجي', 'active'),
('Strategic Impact', 'التأثير الاستراتيجي', 'active');
-- +goose StatementEnd
