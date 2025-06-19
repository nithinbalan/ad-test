-- +goose Up
-- +goose StatementBegin
DECLARE 
  @technicalID INT, @frontendID INT, @backendID INT, @dbID INT, @devopsID INT, 
  @cloudID INT, @apiID INT, @testingID INT, @aiID INT, @mobileID INT,
  @cyberID INT, @threatID INT, @siemID INT, @accessID INT, @dataProtectID INT,
  @complianceKnowID INT, @secureDevID INT, @riskID INT, @controlID INT,
  @registerID INT, @auditID INT, @bcID INT, @thirdPartyID INT, @kpiID INT,
  @projectID INT, @agileID INT, @releaseID INT, @workflowID INT, 
  @raciID INT, @orgLeadID INT, @commID INT, @sectorID INT;

SELECT @technicalID = id FROM competency_groups WHERE title_en = 'Technical Skills';
SELECT @frontendID = id FROM competency_groups WHERE title_en = 'Frontend Development';
SELECT @backendID = id FROM competency_groups WHERE title_en = 'Backend Development';
SELECT @dbID = id FROM competency_groups WHERE title_en = 'Database Management';
SELECT @devopsID = id FROM competency_groups WHERE title_en = 'DevOps & CI/CD';
SELECT @cloudID = id FROM competency_groups WHERE title_en = 'Cloud Infrastructure';
SELECT @apiID = id FROM competency_groups WHERE title_en = 'API Integration';
SELECT @testingID = id FROM competency_groups WHERE title_en = 'Testing & QA';
SELECT @aiID = id FROM competency_groups WHERE title_en = 'AI & Automation';
SELECT @mobileID = id FROM competency_groups WHERE title_en = 'Mobile Development';
SELECT @cyberID = id FROM competency_groups WHERE title_en = 'Cybersecurity & Compliance';
SELECT @threatID = id FROM competency_groups WHERE title_en = 'Threat Management';
SELECT @siemID = id FROM competency_groups WHERE title_en = 'SIEM & Security Tools';
SELECT @accessID = id FROM competency_groups WHERE title_en = 'Access Control & IAM';
SELECT @dataProtectID = id FROM competency_groups WHERE title_en = 'Data Protection';
SELECT @complianceKnowID = id FROM competency_groups WHERE title_en = 'Compliance Knowledge';
SELECT @secureDevID = id FROM competency_groups WHERE title_en = 'Secure Development';
SELECT @riskID = id FROM competency_groups WHERE title_en = 'Risk Management';
SELECT @controlID = id FROM competency_groups WHERE title_en = 'Control Management';
SELECT @registerID = id FROM competency_groups WHERE title_en = 'Compliance Registers';
SELECT @auditID = id FROM competency_groups WHERE title_en = 'Audit Management';
SELECT @bcID = id FROM competency_groups WHERE title_en = 'Business Continuity';
SELECT @thirdPartyID = id FROM competency_groups WHERE title_en = 'Third-Party Risk';
SELECT @kpiID = id FROM competency_groups WHERE title_en = 'KPI/KRI Management';
SELECT @projectID = id FROM competency_groups WHERE title_en = 'Project Planning';
SELECT @agileID = id FROM competency_groups WHERE title_en = 'Agile Methodologies';
SELECT @releaseID = id FROM competency_groups WHERE title_en = 'Release Management';
SELECT @workflowID = id FROM competency_groups WHERE title_en = 'Workflow Automation';
SELECT @raciID = id FROM competency_groups WHERE title_en = 'RACI Matrix Application';
SELECT @orgLeadID = id FROM competency_groups WHERE title_en = 'Organizational & Leadership';
SELECT @commID = id FROM competency_groups WHERE title_en = 'Communication & Documentation';
SELECT @sectorID = id FROM competency_groups WHERE title_en = 'Sector-Specific';

-- Sample inserts
INSERT INTO competencies (competency_group_id, title_en, title_ar, status) VALUES
(@technicalID, 'JavaScript', 'جافا سكريبت', 'active'),
(@technicalID, 'TypeScript', 'تايب سكريبت', 'active'),
(@technicalID, 'Python', 'بايثون', 'active'),
(@technicalID, 'Go', 'جو', 'active'),
(@technicalID, 'C#', 'سي شارب', 'active'),

(@frontendID, 'Angular', 'أنجولار', 'active'),
(@frontendID, 'HTML/CSS', 'إتش تي إم إل / سي إس إس', 'active'),
(@frontendID, 'Responsive UI', 'واجهة مستخدم متجاوبة', 'active'),
(@frontendID, 'Figma', 'فيجما', 'active'),

(@backendID, 'Node.js', 'نود جي إس', 'active'),
(@backendID, '.NET Core', '.نت كور', 'active'),
(@backendID, 'Django', 'دجانغو', 'active'),
(@backendID, 'RESTful APIs', 'واجهات برمجة التطبيقات RESTful', 'active'),

(@dbID, 'Microsoft SQL Server', 'مايكروسوفت إس كيو إل سيرفر', 'active'),
(@dbID, 'MySQL', 'ماي إس كيو إل', 'active'),
(@dbID, 'MongoDB', 'مونغو دي بي', 'active'),

(@devopsID, 'Docker', 'دوكر', 'active'),
(@devopsID, 'Kubernetes', 'كوبرنيتس', 'active'),
(@devopsID, 'Jenkins', 'جينكنز', 'active'),
(@devopsID, 'GitLab CI/CD', 'جيت لاب CI/CD', 'active'),

(@cloudID, 'Azure', 'أزور', 'active'),
(@cloudID, 'AWS', 'خدمات أمازون السحابية', 'active'),
(@cloudID, 'Google Cloud', 'جوجل كلاود', 'active'),
(@cloudID, 'On-Prem Management', 'إدارة محلية', 'active'),

(@apiID, 'REST API', 'واجهة REST', 'active'),
(@apiID, 'SOAP', 'SOAP', 'active'),
(@apiID, 'GraphQL', 'GraphQL', 'active'),

(@testingID, 'Cypress', 'سايبريس', 'active'),
(@testingID, 'Selenium', 'سيلينيوم', 'active'),
(@testingID, 'Postman', 'بوستمان', 'active'),
(@testingID, 'Jest', 'جيست', 'active'),
(@testingID, 'Mocha', 'موكا', 'active'),

(@aiID, 'Prompt Engineering', 'هندسة التعليمات', 'active'),
(@aiID, 'ChatGPT API', 'واجهة ChatGPT', 'active'),
(@aiID, 'RPA Tools', 'أدوات RPA', 'active'),

(@mobileID, 'Flutter', 'فلاتر', 'active'),
(@mobileID, 'Ionic', 'أيونك', 'active'),

(@cyberID, 'NCA ECC', 'المركز الوطني للتصديق الإلكتروني', 'active'),
(@cyberID, 'ISO 27001', 'آيزو 27001', 'active'),
(@cyberID, 'SAMA Cybersecurity', 'الأمن السيبراني SAMA', 'active'),

(@threatID, 'Vulnerability Assessment', 'تقييم الثغرات', 'active'),
(@threatID, 'Threat Modelling', 'نمذجة التهديدات', 'active'),

(@siemID, 'Splunk', 'سبلنك', 'active'),
(@siemID, 'QRadar', 'كيو رادار', 'active'),
(@siemID, 'Tenable', 'تينيبل', 'active'),
(@siemID, 'Nessus', 'نيسوس', 'active'),

(@accessID, 'Active Directory', 'الدليل النشط', 'active'),
(@accessID, 'Role-Based Access Control', 'التحكم بالوصول على أساس الدور', 'active'),

(@dataProtectID, 'Encryption Standards', 'معايير التشفير', 'active'),
(@dataProtectID, 'Data Masking', 'إخفاء البيانات', 'active'),
(@dataProtectID, 'DLP', 'منع تسرب البيانات', 'active'),

(@complianceKnowID, 'GRC Framework Alignment', 'مواءمة إطار الحوكمة والمخاطر والامتثال', 'active'),
(@complianceKnowID, 'Policy Enforcement', 'تطبيق السياسات', 'active'),

(@secureDevID, 'Secure Coding', 'البرمجة الآمنة', 'active'),
(@secureDevID, 'OWASP Top 10', 'أفضل 10 من OWASP', 'active'),

(@riskID, 'Risk Identification', 'تحديد المخاطر', 'active'),
(@riskID, 'Risk Analysis', 'تحليل المخاطر', 'active'),
(@riskID, 'Risk Appetite', 'قابلية المخاطر', 'active'),

(@controlID, 'Control Design', 'تصميم الضوابط', 'active'),
(@controlID, 'Control Testing', 'اختبار الضوابط', 'active'),
(@controlID, 'Efficiency Assessment', 'تقييم الكفاءة', 'active'),

(@registerID, 'Regulatory Mapping', 'رسم خريطة تنظيمية', 'active'),
(@registerID, 'Audit Trails', 'سجلات التدقيق', 'active'),
(@registerID, 'Version Control', 'التحكم في الإصدارات', 'active'),

(@auditID, 'Audit Planning', 'تخطيط التدقيق', 'active'),
(@auditID, 'Findings', 'نتائج التدقيق', 'active'),
(@auditID, 'Corrective Action Workflow', 'سير العمل للإجراءات التصحيحية', 'active'),

(@bcID, 'BIA', 'تحليل الأثر على الأعمال', 'active'),
(@bcID, 'Recovery Planning', 'تخطيط التعافي', 'active'),
(@bcID, 'Crisis Management', 'إدارة الأزمات', 'active'),

(@thirdPartyID, 'Vendor Assessment', 'تقييم البائعين', 'active'),
(@thirdPartyID, 'Contract Risk Evaluation', 'تقييم مخاطر العقود', 'active'),

(@kpiID, 'Threshold Configuration', 'تهيئة العتبة', 'active'),
(@kpiID, 'Performance Monitoring', 'مراقبة الأداء', 'active'),

(@projectID, 'Milestone Management', 'إدارة المعالم', 'active'),
(@projectID, 'Resource Allocation', 'توزيع الموارد', 'active'),

(@agileID, 'Scrum', 'سكرم', 'active'),
(@agileID, 'Kanban', 'كانبان', 'active'),
(@agileID, 'Daily Standups', 'اجتماعات يومية', 'active'),

(@releaseID, 'Versioning', 'إدارة الإصدارات', 'active'),
(@releaseID, 'Change Control', 'ضبط التغيير', 'active'),

(@workflowID, 'BPM', 'إدارة سير العمل', 'active'),
(@workflowID, 'Custom Workflow Design', 'تصميم سير العمل المخصص', 'active'),

(@raciID, 'Role-Activity Mapping', 'ربط الأدوار بالأنشطة', 'active'),
(@raciID, 'Responsibility Allocation', 'توزيع المسؤوليات', 'active'),

(@orgLeadID, 'People Management', 'إدارة الأشخاص', 'active'),
(@orgLeadID, 'Mentoring', 'الإرشاد', 'active'),
(@orgLeadID, 'Cross-functional Communication', 'الاتصال بين الوظائف', 'active'),
(@orgLeadID, 'Feedback Loops', 'حلقات التغذية الراجعة', 'active'),
(@orgLeadID, 'Negotiation Skills', 'مهارات التفاوض', 'active'),
(@orgLeadID, 'Strategic Planning', 'التخطيط الاستراتيجي', 'active'),
(@orgLeadID, 'Knowledge Sharing', 'مشاركة المعرفة', 'active'),

(@commID, 'Formal Email Writing', 'كتابة البريد الإلكتروني الرسمي', 'active'),
(@commID, 'Meeting Notes', 'ملاحظات الاجتماعات', 'active'),
(@commID, 'SOPs', 'إجراءات التشغيل القياسية', 'active'),
(@commID, 'System Architecture Docs', 'مستندات هندسة النظام', 'active'),
(@commID, 'Compliance Reports', 'تقارير الامتثال', 'active'),
(@commID, 'Audit Logs', 'سجلات التدقيق', 'active'),
(@commID, 'Client Demos', 'عروض العملاء', 'active'),
(@commID, 'Internal Reviews', 'مراجعات داخلية', 'active'),

(@sectorID, 'Electrical Network Compliance', 'الامتثال لشبكة الكهرباء', 'active'),
(@sectorID, 'Smart Grid Integration', 'تكامل الشبكة الذكية', 'active'),
(@sectorID, 'NEOM Frameworks', 'أطر نيوم', 'active'),
(@sectorID, 'Public Sector Digitalization', 'التحول الرقمي للقطاع العام', 'active'),
(@sectorID, 'ICS/SCADA Security', 'أمن ICS/SCADA', 'active'),
(@sectorID, 'OT Risk Assessment', 'تقييم مخاطر تكنولوجيا التشغيل', 'active'),
(@sectorID, 'Technical Arabic Translation', 'الترجمة الفنية إلى العربية', 'active'),
(@sectorID, 'Arabic Documentation', 'توثيق باللغة العربية', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE competencies;
-- +goose StatementEnd
