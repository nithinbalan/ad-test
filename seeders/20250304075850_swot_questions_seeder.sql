-- +goose Up
-- +goose StatementBegin
INSERT INTO swot_questions (title_en, title_ar, "order", swot_question_group_id, status) VALUES
    -- Strengths: Resources & Capabilities
    ('What unique skills, assets, or resources do we have?', 'ما هي المهارات أو الأصول أو الموارد الفريدة التي لدينا؟', 1, 1, 'active'),
    ('What do we do better than our competitors?', 'ما الذي نقوم به بشكل أفضل من منافسينا؟', 2, 1, 'active'),
    ('What are our most valuable internal resources?', 'ما هي أهم مواردنا الداخلية قيمة؟', 3, 1, 'active'),

    -- Strengths: Competitive Advantage
    ('What is our unique selling proposition?', 'ما هي ميزتنا التنافسية الفريدة؟', 1, 2, 'active'),
    ('What factors make customers choose us over competitors?', 'ما هي العوامل التي تجعل العملاء يختاروننا على المنافسين؟', 2, 2, 'active'),
    ('Do we have exclusive partnerships or proprietary technology?', 'هل لدينا شراكات حصرية أو تقنيات خاصة؟', 3, 2, 'active'),

    -- Strengths: Market Position
    ('What is our brand reputation in the market?', 'ما هي سمعتنا التجارية في السوق؟', 1, 3, 'active'),
    ('How strong is our customer loyalty?', 'ما مدى قوة ولاء العملاء لدينا؟', 2, 3, 'active'),
    ('Do we have a strong network or industry influence?', 'هل لدينا شبكة قوية أو تأثير في الصناعة؟', 3, 3, 'active'),

    -- Weaknesses: Operational Limitations
    ('What internal processes slow us down?', 'ما هي العمليات الداخلية التي تبطئنا؟', 1, 4, 'active'),
    ('Where are we lacking in efficiency or quality?', 'أين نفتقر إلى الكفاءة أو الجودة؟', 2, 4, 'active'),
    ('Do we have gaps in skills or technology?', 'هل لدينا فجوات في المهارات أو التكنولوجيا؟', 3, 4, 'active'),

    -- Weaknesses: Financial Constraints
    ('Are we financially stable?', 'هل نحن مستقرون ماليًا؟', 1, 5, 'active'),
    ('Do we have cash flow or budget limitations?', 'هل لدينا قيود في التدفق النقدي أو الميزانية؟', 2, 5, 'active'),
    ('Are we overly dependent on a few customers or revenue sources?', 'هل نعتمد بشكل مفرط على عدد قليل من العملاء أو مصادر الإيرادات؟', 3, 5, 'active'),

    -- Weaknesses: Market Challenges
    ('Do we have difficulty in customer acquisition?', 'هل نواجه صعوبة في جذب العملاء؟', 1, 6, 'active'),
    ('Are our marketing and sales efforts effective?', 'هل جهودنا في التسويق والمبيعات فعالة؟', 2, 6, 'active'),
    ('Do we struggle with brand awareness?', 'هل نواجه صعوبة في زيادة الوعي بالعلامة التجارية؟', 3, 6, 'active'),

    -- Opportunities: Market & Industry Trends
    ('What emerging trends can we take advantage of?', 'ما هي الاتجاهات الناشئة التي يمكننا الاستفادة منها؟', 1, 7, 'active'),
    ('Are there new customer needs we can fulfill?', 'هل هناك احتياجات جديدة للعملاء يمكننا تلبيتها؟', 2, 7, 'active'),
    ('Can we expand into new markets or demographics?', 'هل يمكننا التوسع في أسواق أو فئات جديدة؟', 3, 7, 'active'),

    -- Opportunities: Innovation & Technology
    ('Are there new technologies that can enhance our business?', 'هل هناك تقنيات جديدة يمكن أن تحسن أعمالنا؟', 1, 8, 'active'),
    ('Can we improve efficiency through automation?', 'هل يمكننا تحسين الكفاءة من خلال الأتمتة؟', 2, 8, 'active'),
    ('Are there research and development opportunities?', 'هل هناك فرص للبحث والتطوير؟', 3, 8, 'active'),

    -- Opportunities: Strategic Partnerships
    ('Are there potential partnerships or collaborations we can explore?', 'هل هناك شراكات أو تعاونات محتملة يمكننا استكشافها؟', 1, 9, 'active'),
    ('Can we benefit from government incentives or funding?', 'هل يمكننا الاستفادة من الحوافز الحكومية أو التمويل؟', 2, 9, 'active'),
    ('Are there acquisition or investment opportunities?', 'هل هناك فرص للاستحواذ أو الاستثمار؟', 3, 9, 'active'),

    -- Threats: Competition & Market Risks
    ('What threats do our competitors pose?', 'ما هي التهديدات التي يشكلها منافسونا؟', 1, 10, 'active'),
    ('Are there market shifts that could negatively impact us?', 'هل هناك تحولات في السوق قد تؤثر علينا سلبًا؟', 2, 10, 'active'),
    ('Do we face risks from changing regulations?', 'هل نواجه مخاطر من تغير القوانين؟', 3, 10, 'active'),

    -- Threats: Financial & Economic Risks
    ('Are we vulnerable to economic downturns?', 'هل نحن معرضون للركود الاقتصادي؟', 1, 11, 'active'),
    ('Do we have financial dependencies that pose risks?', 'هل لدينا اعتمادات مالية تشكل مخاطر؟', 2, 11, 'active'),
    ('Are we prepared for unexpected expenses or crises?', 'هل نحن مستعدون للنفقات غير المتوقعة أو الأزمات؟', 3, 11, 'active'),

    -- Threats: Operational & Technological Risks
    ('Do we have cybersecurity vulnerabilities?', 'هل لدينا نقاط ضعف في الأمن السيبراني؟', 1, 12, 'active'),
    ('Are we at risk of supply chain disruptions?', 'هل نحن معرضون لخطر تعطل سلسلة التوريد؟', 2, 12, 'active'),
    ('Do we have contingency plans for emergencies?', 'هل لدينا خطط طوارئ لمواجهة الأزمات؟', 3, 12, 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE swot_questions;
-- +goose StatementEnd
