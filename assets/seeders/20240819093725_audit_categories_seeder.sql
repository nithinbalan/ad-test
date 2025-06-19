-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_categories
    (title_en, title_ar, description_en, description_ar, status)
VALUES
    ('Financial Audit', 'التدقيق المالي', 'An examination of an organization''s financial statements and related processes to ensure accuracy and compliance with accounting standards.', 'فحص للبيانات المالية للمنظمة والعمليات ذات الصلة لضمان الدقة والامتثال لمعايير المحاسبة.', 'active'),
    ('Operational Audit', 'التدقيق التشغيلي', 'A review of an organization''s operations to evaluate efficiency, effectiveness, and adherence to management policies.', 'مراجعة عمليات المنظمة لتقييم الكفاءة والفعالية والالتزام بسياسات الإدارة.', 'active'),
    ('Compliance Audit', 'تدقيق الامتثال', 'An audit to ensure that an organization adheres to external laws, regulations, and internal policies.', 'تدقيق لضمان أن المنظمة تلتزم بالقوانين واللوائح الخارجية والسياسات الداخلية.', 'active'),
    ('Information Technology (IT) Audit', 'تدقيق تكنولوجيا المعلومات', 'An evaluation of an organization''s IT infrastructure, applications, and processes to ensure data integrity and security.', 'تقييم للبنية التحتية لتكنولوجيا المعلومات والتطبيقات والعمليات في المنظمة لضمان سلامة وأمان البيانات.', 'active'),
    ('Environmental Audit', 'التدقيق البيئي', 'An assessment of an organization''s environmental management practices and their impact on the environment.', 'تقييم لممارسات إدارة البيئة في المنظمة وتأثيرها على البيئة.', 'active'),
    ('Performance Audit', 'تدقيق الأداء', 'A review focused on assessing whether resources are used efficiently and effectively to achieve organizational goals.', 'مراجعة تركز على تقييم ما إذا كانت الموارد تُستخدم بكفاءة وفعالية لتحقيق أهداف المنظمة.', 'active'),
    ('Internal Audit', 'التدقيق الداخلي', 'An independent, objective evaluation of an organization''s internal controls, processes, and risk management practices.', 'تقييم مستقل وموضوعي لضوابط وعمليات وإدارة المخاطر الداخلية في المنظمة.', 'active'),
    ('External Audit', 'التدقيق الخارجي', 'An independent examination of an organization''s financial statements conducted by an external auditor.', 'فحص مستقل للبيانات المالية للمنظمة يجريه مدقق خارجي.', 'active'),
    ('Risk-Based Audit', 'التدقيق المعتمد على المخاطر', 'An audit focused on identifying and evaluating the risks that could affect an organization''s objectives.', 'تدقيق يركز على تحديد وتقييم المخاطر التي قد تؤثر على أهداف المنظمة.', 'active'),
    ('Follow-Up Audit', 'تدقيق المتابعة', 'A review to ensure that corrective actions from previous audits have been implemented effectively.', 'مراجعة لضمان تنفيذ الإجراءات التصحيحية من التدقيقات السابقة بفعالية.', 'active'),
    ('Special Investigation Audit', 'تدقيق التحقيق الخاص', 'An audit conducted to investigate specific issues, such as fraud, misconduct, or irregularities.', 'تدقيق يجرى للتحقيق في قضايا معينة مثل الاحتيال أو سوء السلوك أو المخالفات.', 'active'),
    ('Quality Management Audit', 'تدقيق إدارة الجودة', 'An evaluation of an organization''s quality management system to ensure compliance with quality standards.', 'تقييم لنظام إدارة الجودة في المنظمة لضمان الامتثال لمعايير الجودة.', 'active'),
    ('Health and Safety Audit', 'تدقيق الصحة والسلامة', 'A review of an organization''s health and safety practices to ensure a safe working environment.', 'مراجعة لممارسات الصحة والسلامة في المنظمة لضمان بيئة عمل آمنة.', 'active'),
    ('Process Audit', 'تدقيق العمليات', 'An audit focused on evaluating the efficiency and effectiveness of a specific process within the organization.', 'تدقيق يركز على تقييم كفاءة وفعالية عملية محددة داخل المنظمة.', 'active'),
    ('Forensic Audit', 'التدقيق الجنائي', 'An audit focused on investigating financial crimes such as fraud and embezzlement.', 'تدقيق يركز على التحقيق في الجرائم المالية مثل الاحتيال والاختلاس.', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_categories;
-- +goose StatementEnd
