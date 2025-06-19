-- +goose Up
-- +goose StatementBegin

UPDATE swot_pestel_questions
SET is_swot = 1, is_pestel = 0
WHERE swot_pestel_question_group_id IN (1, 2, 3, 4, 5, 6, 7, 8,9,10,11,12,13,14,15,16);

INSERT INTO swot_pestel_questions (title_en, title_ar, "order", swot_pestel_question_group_id, status, is_swot, is_pestel)
VALUES
    -- Political Factors (PESTEL)
    ('What government policies impact our industry?', 'ما هي السياسات الحكومية التي تؤثر على صناعتنا؟', 1, 17, 'active', 0, 1),
    ('Are there upcoming regulatory changes?', 'هل هناك تغييرات تنظيمية قادمة؟', 2, 17, 'active', 0, 1),
    ('How stable is the political environment?', 'ما مدى استقرار البيئة السياسية؟', 3, 17, 'active', 0, 1),

    -- Economic Factors (PESTEL)
    ('How do current economic trends affect us?', 'كيف تؤثر الاتجاهات الاقتصادية الحالية علينا؟', 1, 18, 'active', 0, 1),
    ('What is the impact of inflation on our costs?', 'ما هو تأثير التضخم على تكاليفنا؟', 2, 18, 'active', 0, 1),
    ('Are interest rates affecting our business?', 'هل تؤثر أسعار الفائدة على أعمالنا؟', 3, 18, 'active', 0, 1),

    -- Social Factors (PESTEL)
    ('Are there demographic shifts affecting our market?', 'هل هناك تغييرات ديموغرافية تؤثر على سوقنا؟', 1, 19, 'active', 0, 1),
    ('How do cultural trends influence our customers?', 'كيف تؤثر الاتجاهات الثقافية على عملائنا؟', 2, 19, 'active', 0, 1),
    ('What are the key societal values in our industry?', 'ما هي القيم المجتمعية الرئيسية في صناعتنا؟', 3, 19, 'active', 0, 1),

    -- Technological Factors (PESTEL)
    ('What emerging technologies could disrupt our market?', 'ما هي التقنيات الناشئة التي قد تعطل سوقنا؟', 1, 20, 'active', 0, 1),
    ('Are we keeping up with technological advancements?', 'هل نواكب التقدم التكنولوجي؟', 2, 20, 'active', 0, 1),
    ('What R&D investments can give us a competitive edge?', 'ما هي الاستثمارات في البحث والتطوير التي يمكن أن تعطينا ميزة تنافسية؟', 3, 20, 'active', 0, 1),

    -- Environmental Factors (PESTEL)
    ('Are we compliant with environmental regulations?', 'هل نحن ممتثلون للوائح البيئية؟', 1, 21, 'active', 0, 1),
    ('How do climate changes affect our business?', 'كيف تؤثر التغيرات المناخية على أعمالنا؟', 2, 21, 'active', 0, 1),
    ('What sustainability initiatives can we adopt?', 'ما هي المبادرات البيئية التي يمكننا تبنيها؟', 3, 21, 'active', 0, 1),

    -- Legal Factors (PESTEL)
    ('Are there new industry laws affecting us?', 'هل هناك قوانين صناعية جديدة تؤثر علينا؟', 1, 22, 'active', 0, 1),
    ('How do labor laws impact our workforce?', 'كيف تؤثر قوانين العمل على قوتنا العاملة؟', 2, 22, 'active', 0, 1),
    ('What are the legal risks in our sector?', 'ما هي المخاطر القانونية في قطاعنا؟', 3, 22, 'active', 0, 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE swot_pestel_questions;
-- +goose StatementEnd
