-- +goose Up
-- +goose StatementBegin

DECLARE @ComplianceRiskID INT, @FinancialRiskID INT, @OperationalRiskID INT, @ReputationalRiskID INT, @StrategicRiskID INT;

-- Fetching Risk Category IDs
SELECT @ComplianceRiskID = id FROM risk_categories WHERE title_en = 'Compliance Risk';
SELECT @FinancialRiskID = id FROM risk_categories WHERE title_en = 'Financial Risk';
SELECT @OperationalRiskID = id FROM risk_categories WHERE title_en = 'Operational Risk';
SELECT @ReputationalRiskID = id FROM risk_categories WHERE title_en = 'Reputational Risk';
SELECT @StrategicRiskID = id FROM risk_categories WHERE title_en = 'Strategic Risk';

-- Inserting Risk Appetite Statements
INSERT INTO risk_appetite_statements (reference_number, risk_category_id, statement_en, statement_ar, status, created_by_id, created_at)
VALUES
    -- Compliance Risk Statements
    ('RAS-2024-0001', @ComplianceRiskID, 
     'Our organization maintains a zero-tolerance policy towards non-compliance with regulatory requirements. We are committed to meeting all local and international regulations to avoid any legal penalties.',
     'تحافظ منظمتنا على سياسة عدم التسامح مع عدم الامتثال للمتطلبات التنظيمية. نحن ملتزمون بالامتثال لجميع اللوائح المحلية والدولية لتجنب أي عقوبات قانونية.',
     'active', 1, GETDATE()),

    ('RAS-2024-0002', @ComplianceRiskID, 
     'We have a low risk appetite for breaches of internal policies. All employees are expected to adhere to internal guidelines, and any deviations will be promptly addressed.',
     'لدينا شهية منخفضة تجاه خرق السياسات الداخلية. من المتوقع أن يلتزم جميع الموظفين بالإرشادات الداخلية، وسيتم معالجة أي انحرافات على الفور.',
     'active', 1, GETDATE()),

    ('RAS-2024-0003', @ComplianceRiskID, 
     'Our organization expects all third-party partners to comply with our compliance standards. We have a moderate risk appetite for potential third-party non-compliance, but we require corrective actions to be taken immediately.',
     'تتوقع منظمتنا أن يمتثل جميع الشركاء الخارجيين لمعايير الامتثال لدينا. لدينا شهية معتدلة تجاه عدم الامتثال المحتمل من قبل الأطراف الثالثة، ولكننا نطلب اتخاذ إجراءات تصحيحية على الفور.',
     'active', 1, GETDATE()),

    -- Financial Risk Statements
    ('RAS-2024-0004', @FinancialRiskID, 
     'We maintain a low risk appetite for budget overruns and prioritize strict financial controls to ensure that expenditures remain within approved limits.',
     'نحافظ على شهية منخفضة تجاه تجاوز الميزانية ونولي الأولوية للضوابط المالية الصارمة لضمان بقاء النفقات ضمن الحدود المعتمدة.',
     'active', 1, GETDATE()),

    ('RAS-2024-0005', @FinancialRiskID, 
     'Our organization has a moderate risk appetite for investments in new ventures and technologies. We are willing to accept calculated financial risks that align with our strategic growth objectives.',
     'لدينا شهية معتدلة تجاه الاستثمارات في المشاريع الجديدة والتكنولوجيات. نحن على استعداد لقبول المخاطر المالية المحسوبة التي تتماشى مع أهدافنا الاستراتيجية للنمو.',
     'active', 1, GETDATE()),

    ('RAS-2024-0006', @FinancialRiskID, 
     'We accept a moderate level of risk in revenue variability, recognizing that fluctuations may occur due to market conditions. We aim to mitigate this through diversified income streams.',
     'نقبل مستوى معتدلًا من المخاطر في تباين الإيرادات، معترفًا بأن التقلبات قد تحدث بسبب ظروف السوق. نحن نسعى للتخفيف من ذلك من خلال تنويع مصادر الدخل.',
     'active', 1, GETDATE()),

    -- Operational Risk Statements
    ('RAS-2024-0007', @OperationalRiskID, 
     'Our organization accepts a moderate level of operational risk in pursuit of process efficiency. We invest in continuous improvement and innovation to enhance our operational effectiveness.',
     'تقبل منظمتنا مستوى معتدلًا من المخاطر التشغيلية في السعي لتحقيق كفاءة العمليات. نحن نستثمر في التحسين المستمر والابتكار لتعزيز فعاليتنا التشغيلية.',
     'active', 1, GETDATE()),

    ('RAS-2024-0008', @OperationalRiskID, 
     'We have a moderate to high risk appetite for the adoption of new technologies, understanding that some operational disruption may occur during implementation. We manage this risk through careful planning and training.',
     'لدينا شهية معتدلة إلى عالية تجاه اعتماد التكنولوجيات الجديدة، مع فهم أن بعض الاضطرابات التشغيلية قد تحدث أثناء التنفيذ. نحن ندير هذه المخاطر من خلال التخطيط الدقيق والتدريب.',
     'active', 1, GETDATE()),

    ('RAS-2024-0009', @OperationalRiskID, 
     'Our organization maintains a low risk appetite for disruptions in the supply chain. We prioritize robust supplier relationships and contingency planning to ensure continuity of operations.',
     'نحافظ على شهية منخفضة تجاه الاضطرابات في سلسلة التوريد. نحن نولي الأولوية للعلاقات القوية مع الموردين والتخطيط للطوارئ لضمان استمرارية العمليات.',
     'active', 1, GETDATE()),

    -- Reputational Risk Statements
    ('RAS-2024-0010', @ReputationalRiskID, 
     'We have a low risk appetite for actions that could harm our brand image. All initiatives are carefully reviewed to ensure they align with our values and public perception.',
     'لدينا شهية منخفضة تجاه الإجراءات التي قد تضر بصورة علامتنا التجارية. تتم مراجعة جميع المبادرات بعناية للتأكد من توافقها مع قيمنا والتصور العام.',
     'active', 1, GETDATE()),

    ('RAS-2024-0011', @ReputationalRiskID, 
     'Our organization is cautious with media exposure, maintaining a low risk appetite for any public communications that could be misinterpreted or damage our reputation.',
     'نحن حذرون من التعرض الإعلامي، حيث نحافظ على شهية منخفضة تجاه أي اتصالات عامة يمكن أن تُساء تفسيرها أو تضر بسمعتنا.',
     'active', 1, GETDATE()),

    ('RAS-2024-0012', @ReputationalRiskID, 
     'We prioritize maintaining customer trust and have a very low risk appetite for any actions that could erode this trust. Customer satisfaction is central to our reputation management strategy.',
     'نحن نولي الأولوية للحفاظ على ثقة العملاء ولدينا شهية منخفضة للغاية تجاه أي إجراءات قد تؤدي إلى تآكل هذه الثقة. رضا العملاء هو محور استراتيجيتنا لإدارة السمعة.',
     'active', 1, GETDATE()),

    -- Strategic Risk Statements
    ('RAS-2024-0013', @StrategicRiskID, 
     'Our organization has a moderate to high risk appetite for strategic risks related to market expansion. We are open to entering new markets and industries where there is potential for significant growth.',
     'لدى منظمتنا شهية معتدلة إلى عالية تجاه المخاطر الاستراتيجية المتعلقة بالتوسع في الأسواق. نحن منفتحون على دخول أسواق وصناعات جديدة حيث توجد إمكانيات لنمو كبير.',
     'active', 1, GETDATE()),

    ('RAS-2024-0014', @StrategicRiskID, 
     'We accept a high level of risk in innovation and research and development (R&D) activities. Our strategy is to invest in cutting-edge technologies and products to maintain competitive advantage.',
     'نحن نقبل مستوى عالٍ من المخاطر في أنشطة الابتكار والبحث والتطوير. استراتيجيتنا هي الاستثمار في التكنولوجيات والمنتجات المتقدمة للحفاظ على الميزة التنافسية.',
     'active', 1, GETDATE()),

    ('RAS-2024-0015', @StrategicRiskID, 
     'We have a moderate risk appetite for forming new partnerships and alliances, recognizing the potential for both significant benefits and challenges. We approach these opportunities with due diligence and strategic alignment.',
     'لدينا شهية معتدلة لتشكيل شراكات وتحالفات جديدة، معترفًا بالإمكانيات الكبيرة للفوائد والتحديات. نحن نتعامل مع هذه الفرص بعناية واجتهاد وتوافق استراتيجي.',
     'active', 1, GETDATE());

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_appetite_statements;
-- +goose StatementEnd
