-- +goose Up
-- +goose StatementBegin
INSERT INTO risk_sources (title_en, title_ar, status) VALUES
('Other', 'آخر', 'active'),
('Human Resources', 'الموارد البشرية', 'active'),
('Technology', 'التكنولوجيا', 'active'),
('Financial Management', 'الإدارة المالية', 'active'),
('Supply Chain', 'سلسلة التوريد', 'active'),
('Market Conditions', 'ظروف السوق', 'active'),
('Natural Disasters', 'الكوارث الطبيعية', 'active'),
('Reputational', 'السمعة', 'active'),
('Strategic Decisions', 'القرارات الاستراتيجية', 'active'),
('Third-Party Relationships', 'العلاقات مع الأطراف الثالثة', 'active'),
('Political and Economic', 'السياسية والاقتصادية', 'active'),
('Environmental', 'البيئة', 'active'),
('Health and Safety', 'الصحة والسلامة', 'active'),
('Compliance', 'الامتثال', 'active'),
('Innovation and R&D', 'الابتكار والبحث والتطوير', 'active'),
('Customer Relations', 'علاقات العملاء', 'active'),
('Data Privacy', 'خصوصية البيانات', 'active'),
('Cultural and Social', 'الثقافية والاجتماعية', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_sources;
-- +goose StatementEnd
