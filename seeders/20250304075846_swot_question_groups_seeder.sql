-- +goose Up
-- +goose StatementBegin
INSERT INTO swot_pestel_question_groups (title_en, title_ar, "order", status) VALUES
    ('Resources & Capabilities', 'الموارد والقدرات', 1, 'active'),
    ('Competitive Advantage', 'الميزة التنافسية', 2, 'active'),
    ('Customer Perception & Reputation', 'تصور العملاء والسمعة', 3, 'active'),
    ('Financial Stability', 'الاستقرار المالي', 4, 'active'),
    ('Operational Limitations', 'القيود التشغيلية', 5, 'active'),
    ('Skill Gaps & Team Challenges', 'فجوات المهارات وتحديات الفريق', 6, 'active'),
    ('Customer & Market Perception Issues', 'مشاكل تصور العملاء والسوق', 7, 'active'),
    ('Financial Constraints', 'القيود المالية', 8, 'active'),
    ('Industry Trends & Market Expansion', 'اتجاهات الصناعة وتوسع السوق', 9, 'active'),
    ('Technology & Innovation', 'التكنولوجيا والابتكار', 10, 'active'),
    ('Strategic Partnerships & Collaborations', 'الشراكات الاستراتيجية والتعاون', 11, 'active'),
    ('Regulatory & Economic Opportunities', 'الفرص التنظيمية والاقتصادية', 12, 'active'),
    ('Competition & Market Disruption', 'المنافسة والاضطرابات السوقية', 13, 'active'),
    ('Economic & Financial Risks', 'المخاطر الاقتصادية والمالية', 14, 'active'),
    ('Technological & Industry Changes', 'التغيرات التكنولوجية والصناعية', 15, 'active'),
    ('Regulatory & Supply Chain Risks', 'المخاطر التنظيمية وسلسلة التوريد', 16, 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE swot_pestel_question_groups;
-- +goose StatementEnd
