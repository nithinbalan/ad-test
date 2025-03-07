-- +goose Up
-- +goose StatementBegin
UPDATE swot_pestel_question_groups
SET is_swot = 1, is_pestel = 0
WHERE title_en IN (
    'Resources & Capabilities',
    'Competitive Advantage',
    'Customer Perception & Reputation',
    'Financial Stability',
    'Operational Limitations',
    'Skill Gaps & Team Challenges',
    'Customer & Market Perception Issues',
    'Financial Constraints',
    'Industry Trends & Market Expansion',
    'Technology & Innovation',
    'Strategic Partnerships & Collaborations',                  
    'Regulatory & Economic Opportunities', 
    'Competition & Market Disruption', 
    'Economic & Financial Risks', 
    'Technological & Industry Changes', 
    'Regulatory & Supply Chain Risks'
);

INSERT INTO swot_pestel_question_groups (title_en, title_ar, "order", status, is_swot, is_pestel)
VALUES
    ('Political Factors', 'العوامل السياسية', 1, 'active', 0, 1),
    ('Economic Factors', 'العوامل الاقتصادية', 2, 'active', 0, 1),
    ('Social Factors', 'العوامل الاجتماعية', 3, 'active', 0, 1),
    ('Technological Factors', 'العوامل التكنولوجية', 4, 'active', 0, 1),
    ('Environmental Factors', 'العوامل البيئية', 5, 'active', 0, 1),
    ('Legal Factors', 'العوامل القانونية', 6, 'active', 0, 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE swot_pestel_question_groups;
-- +goose StatementEnd
