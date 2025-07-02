-- +goose Up
-- +goose StatementBegin

DECLARE @productionStatusId INT, @comingSoonStatusId INT;

SELECT @productionStatusId = id FROM knowledge_center_resource_statuses WHERE title_en = 'Production';
SELECT @comingSoonStatusId = id FROM knowledge_center_resource_statuses WHERE title_en = 'Coming Soon';

DECLARE @moduleID INT;

-- Risk Management
INSERT INTO ai_modules (title_en, title_ar, [order]) VALUES (N'Risk Management', N'إدارة المخاطر', 1);
SELECT @moduleID = SCOPE_IDENTITY();
INSERT INTO ai_module_features (ai_module_id, title_en, title_ar, knowledge_center_resource_status_id, [order]) VALUES
(@moduleID, N'Risk analysis recommendations based on similar incidents, impact, and likelihood', N'توصيات تحليل المخاطر استنادًا إلى حوادث مماثلة والتأثير والاحتمالية', @comingSoonStatusId, 1),
(@moduleID, N'Risk mitigation strategy suggestions', N'اقتراحات لاستراتيجيات تخفيف المخاطر', @comingSoonStatusId, 2),
(@moduleID, N'Auto-suggest risk categories and related controls', N'اقتراح تلقائي لفئات المخاطر والضوابط ذات الصلة', @comingSoonStatusId, 3),
(@moduleID, N'Risk rating estimation based on ISO 31000 / COSO', N'تقدير تقييم المخاطر بناءً على ISO 31000 / COSO', @productionStatusId, 4),
(@moduleID, N'Detect duplicates or redundant risks across units', N'الكشف عن التكرارات أو المخاطر الزائدة عبر الوحدات', @comingSoonStatusId, 5);

-- Audit Management
INSERT INTO ai_modules (title_en, title_ar, [order]) VALUES (N'Audit Management', N'إدارة التدقيق', 2);
SELECT @moduleID = SCOPE_IDENTITY();
INSERT INTO ai_module_features (ai_module_id, title_en, title_ar, knowledge_center_resource_status_id, [order]) VALUES
(@moduleID, N'AI-suggested audit findings based on transactional and checklist data', N'نتائج تدقيق مقترحة من الذكاء الاصطناعي استنادًا إلى البيانات والمعايير', @productionStatusId, 1),
(@moduleID, N'Pattern-based fraud detection and anomaly flagging', N'كشف الاحتيال بناءً على الأنماط وتحديد الشذوذ', @productionStatusId, 2),
(@moduleID, N'Automated recommendations for corrective actions', N'توصيات تلقائية للإجراءات التصحيحية', @productionStatusId, 3),
(@moduleID, N'Draft and final report review with keyword-based insight extraction', N'مراجعة التقارير مع استخراج رؤى بناءً على الكلمات المفتاحية', @comingSoonStatusId, 4),
(@moduleID, N'Identify gaps between planned and performed audits', N'تحديد الفجوات بين التدقيق المخطط والمنفذ', @comingSoonStatusId, 5),
(@moduleID, N'Risk-focused observation suggestions based on department-wise trends', N'ملاحظات مركزة على المخاطر استنادًا إلى الاتجاهات حسب القسم', @comingSoonStatusId, 6),
(@moduleID, N'Action plan quality scoring and duplicate detection', N'تقييم جودة خطة العمل وكشف التكرار', @comingSoonStatusId, 7);

-- Compliance Management
INSERT INTO ai_modules (title_en, title_ar, [order]) VALUES (N'Compliance Management', N'إدارة الامتثال', 3);
SELECT @moduleID = SCOPE_IDENTITY();
INSERT INTO ai_module_features (ai_module_id, title_en, title_ar, knowledge_center_resource_status_id, [order]) VALUES
(@moduleID, N'AI-powered mapping of policies to applicable compliance requirements (e.g., ISO, NCA, SAMA)', N'ربط السياسات بمتطلبات الامتثال', @comingSoonStatusId, 1),
(@moduleID, N'Automatic association of policies with controls, processes, and risks', N'الربط التلقائي للسياسات مع الضوابط والعمليات والمخاطر', @comingSoonStatusId, 2),
(@moduleID, N'Identification of gaps or missing compliance clauses', N'تحديد الفجوات أو بنود الامتثال المفقودة', @comingSoonStatusId, 3),
(@moduleID, N'Suggest improvements for regulatory alignment', N'اقتراح تحسينات للتوافق التنظيمي', @comingSoonStatusId, 4),
(@moduleID, N'Generate compliance summaries per standard', N'إنشاء ملخصات الامتثال حسب المعايير', @comingSoonStatusId, 5),
(@moduleID, N'Flag outdated or conflicting compliance references in documents', N'وضع علامات على المراجع القديمة أو المتضاربة', @productionStatusId, 6);

-- Policy & Document Intelligence
INSERT INTO ai_modules (title_en, title_ar, [order]) VALUES (N'Policy & Document Intelligence', N'ذكاء السياسات والوثائق', 4);
SELECT @moduleID = SCOPE_IDENTITY();
INSERT INTO ai_module_features (ai_module_id, title_en, title_ar, knowledge_center_resource_status_id, [order]) VALUES
(@moduleID, N'Summarize lengthy documents into digestible key points', N'تلخيص المستندات الطويلة إلى نقاط رئيسية', @comingSoonStatusId, 1),
(@moduleID, N'Identify missing or conflicting clauses in policies', N'تحديد البنود المفقودة أو المتضاربة في السياسات', @comingSoonStatusId, 2),
(@moduleID, N'Provide clause-level extraction and classification', N'استخراج البنود وتصنيفها', @comingSoonStatusId, 3),
(@moduleID, N'Enable natural language Q&A over policy content', N'دعم الأسئلة والأجوبة بلغة طبيعية حول محتوى السياسات', @comingSoonStatusId, 4),
(@moduleID, N'Detect outdated terms, references, or regulatory misalignments', N'كشف المصطلحات القديمة أو الانحرافات التنظيمية', @comingSoonStatusId, 5),
(@moduleID, N'Auto-tag and classify documents based on context', N'وضع علامات وتصنيف المستندات تلقائيًا', @comingSoonStatusId, 6),
(@moduleID, N'Assist in preparing or updating policies with AI-generated suggestions', N'المساعدة في إعداد أو تحديث السياسات باقتراحات الذكاء الاصطناعي', @productionStatusId, 7);

-- Conversational Assistant
INSERT INTO ai_modules (title_en, title_ar, [order]) VALUES (N'Conversational Assistant (Beta)', N'المساعد التفاعلي (تجريبي)', 5);
SELECT @moduleID = SCOPE_IDENTITY();
INSERT INTO ai_module_features (ai_module_id, title_en, title_ar, knowledge_center_resource_status_id, [order]) VALUES
(@moduleID, N'Interact with GRC data using plain language queries', N'التفاعل مع بيانات GRC باستخدام الاستعلامات الطبيعية', @comingSoonStatusId, 1),
(@moduleID, N'Context-aware responses drawn from multiple modules', N'ردود قائمة على السياق من عدة وحدات', @comingSoonStatusId, 2),
(@moduleID, N'Natural language support for multilingual environments', N'دعم اللغة الطبيعية لبيئات متعددة اللغات', @comingSoonStatusId, 3),
(@moduleID, N'Voice-to-text integration for hands-free queries', N'تكامل الصوت إلى نص للاستعلامات بدون استخدام اليدين', @comingSoonStatusId, 4),
(@moduleID, N'Provide explanations of terms, frameworks, and statuses', N'تقديم تفسيرات للمصطلحات والأطر والحالات', @comingSoonStatusId, 5),
(@moduleID, N'Route user to relevant section upon AI query', N'توجيه المستخدم إلى القسم المناسب عند الاستعلام', @comingSoonStatusId, 6);

-- Asset Management
INSERT INTO ai_modules (title_en, title_ar, [order]) VALUES (N'Asset Management', N'إدارة الأصول', 6);
SELECT @moduleID = SCOPE_IDENTITY();
INSERT INTO ai_module_features (ai_module_id, title_en, title_ar, knowledge_center_resource_status_id, [order]) VALUES
(@moduleID, N'Recommend asset criticality score based on impact and dependencies', N'توصية بدرجة أهمية الأصل بناءً على التأثير والاعتمادية', @productionStatusId, 1),
(@moduleID, N'Identify key vulnerabilities and related threats linked to assets', N'تحديد الثغرات والتهديدات المرتبطة بالأصول', @comingSoonStatusId, 2),
(@moduleID, N'Auto-suggest controls based on asset type and risk exposure', N'اقتراح الضوابط بناءً على نوع الأصل والتعرض للخطر', @comingSoonStatusId, 3),
(@moduleID, N'Map assets with policies, risks, and processes using AI', N'ربط الأصول بالسياسات والمخاطر والعمليات باستخدام الذكاء الاصطناعي', @comingSoonStatusId, 4),
(@moduleID, N'Generate alerts for high-risk assets without mitigation', N'إنشاء تنبيهات للأصول عالية المخاطر بدون تخفيف', @comingSoonStatusId, 5);

-- Business Continuity Management
INSERT INTO ai_modules (title_en, title_ar, [order]) VALUES (N'Business Continuity Management (BCM)', N'إدارة استمرارية الأعمال', 7);
SELECT @moduleID = SCOPE_IDENTITY();
INSERT INTO ai_module_features (ai_module_id, title_en, title_ar, knowledge_center_resource_status_id, [order]) VALUES
(@moduleID, N'Suggest RTO and RPO', N'اقتراح RTO و RPO', @comingSoonStatusId, 1),
(@moduleID, N'Identify critical business functions and interdependencies', N'تحديد الوظائف الحيوية والاعتماديات', @comingSoonStatusId, 2),
(@moduleID, N'Recommend continuity strategies based on past disruptions', N'اقتراح استراتيجيات الاستمرارية استنادًا إلى الانقطاعات السابقة', @comingSoonStatusId, 3),
(@moduleID, N'Provide AI-generated incident response checklists', N'توفير قوائم تحقق للاستجابة للحوادث من الذكاء الاصطناعي', @comingSoonStatusId, 4),
(@moduleID, N'Highlight gaps in plans by comparing with standards', N'تسليط الضوء على الفجوات في الخطط عبر المقارنة بالمعايير', @comingSoonStatusId, 5);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE ai_module_features;
TRUNCATE TABLE ai_modules;
-- +goose StatementEnd
