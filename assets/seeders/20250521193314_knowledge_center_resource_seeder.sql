-- +goose Up
-- +goose StatementBegin

-- Declare module ID variables
DECLARE @complianceID INT, @riskID INT, @documentID INT;
DECLARE @auditID INT, @strategyID INT, @assetID INT, @knowledgeID INT;

-- Declare status ID variable
DECLARE @newStatusID INT, @comingSoonStatusID INT, @productionStatusID INT;

-- Assign values from modules table
SELECT @complianceID = id FROM modules WHERE identifier = 'compliance';
SELECT @riskID = id FROM modules WHERE identifier = 'risk';
SELECT @documentID = id FROM modules WHERE identifier = 'document';
SELECT @auditID = id FROM modules WHERE identifier = 'audit';
SELECT @strategyID = id FROM modules WHERE identifier = 'strategy';
SELECT @assetID = id FROM modules WHERE identifier = 'asset';
SELECT @knowledgeID = id FROM modules WHERE identifier = 'knowledge_center';

-- Get knowledge_center_resource_status_id 
SELECT @newStatusID = id FROM knowledge_center_resource_statuses WHERE title_en = 'New';
SELECT @productionStatusID = id FROM knowledge_center_resource_statuses WHERE title_en = 'Production';
SELECT @comingSoonStatusID = id FROM knowledge_center_resource_statuses WHERE title_en = 'Coming Soon';


-- Compliance Management
INSERT INTO knowledge_center_resources (module_id, knowledge_center_resource_status_id, title_en, title_ar, icon, created_at, updated_at) VALUES
(@complianceID, @comingSoonStatusID, 'Saudi Arabian Monetary Authority (SAMA)', 'مؤسسة النقد العربي السعودي (ساما)', 'sama-icon.png', GETDATE(), NULL),
(@complianceID, @comingSoonStatusID, 'National Cybersecurity Authority (NCA)', '(NCA) الهيئة الوطنية للأمن السيبراني', 'nca-icon.png', GETDATE(), NULL),
(@complianceID, @comingSoonStatusID, 'ISO/IEC 27001', 'آيزو/آي إي سي 27001', 'iso-icon.png', GETDATE(), NULL);

-- Risk Management
INSERT INTO knowledge_center_resources (module_id, knowledge_center_resource_status_id, title_en, title_ar, icon, created_at, updated_at) VALUES
(@riskID, @newStatusID, 'Risk Library', 'مكتبة المخاطر', 'risk-library-icon.png', GETDATE(), NULL),
(@riskID, @newStatusID, 'Key Risk Indicators Library', 'مكتبة مؤشرات المخاطر الرئيسية', 'kri-icon.png', GETDATE(), NULL),
(@riskID, @comingSoonStatusID, 'KPI Library', 'مكتبة مؤشرات الأداء الرئيسية', 'kpi-icon.png', GETDATE(), NULL);

-- Document Management
INSERT INTO knowledge_center_resources (module_id, knowledge_center_resource_status_id, title_en, title_ar, icon, created_at, updated_at) VALUES
(@documentID, @comingSoonStatusID, 'Policy Templates Library', 'مكتبة قوالب السياسة', 'policy-library-icon.png', GETDATE(), NULL),
(@documentID, @comingSoonStatusID, 'Standards Library', 'مكتبة المعايير', 'standards-library-icon.png', GETDATE(), NULL);

-- Audit Management
INSERT INTO knowledge_center_resources (module_id, knowledge_center_resource_status_id, title_en, title_ar, icon, created_at, updated_at) VALUES
(@auditID, @comingSoonStatusID, 'Audit Checklist Library', 'مكتبة قوائم التدقيق', 'audit-library-icon.png', GETDATE(), NULL);

-- Strategy Management
INSERT INTO knowledge_center_resources (module_id, knowledge_center_resource_status_id, title_en, title_ar, icon, created_at, updated_at) VALUES
(@strategyID, @comingSoonStatusID, 'Strategic Objectives Library', 'مكتبة الأهداف الاستراتيجية', 'strategic-library-icon.png', GETDATE(), NULL);

-- Asset Management
INSERT INTO knowledge_center_resources (module_id, knowledge_center_resource_status_id, title_en, title_ar, icon, created_at, updated_at) VALUES
(@assetID, @comingSoonStatusID, 'Asset Library', 'مكتبة الأصول', 'asset-library-icon.png', GETDATE(), NULL),
(@assetID, @comingSoonStatusID, 'Vulnerabilities Library', 'مكتبة الثغرات الأمنية', 'vulnerabilities-library-icon.png', GETDATE(), NULL),
(@assetID, @comingSoonStatusID, 'Threats Library', 'مكتبة التهديدات', 'threats-library-icon.png', GETDATE(), NULL);

-- Knowledge Resources
INSERT INTO knowledge_center_resources (module_id, knowledge_center_resource_status_id, title_en, title_ar, icon, created_at, updated_at) VALUES
(@knowledgeID, @productionStatusID, 'Cyber Security', 'الأمن السيبراني', 'cyber-security-icon.png', GETDATE(), NULL),
(@knowledgeID, @productionStatusID, 'AI & Generative AI', 'الذكاء الاصطناعي و الذكاء التوليدي', 'ai-generative-ai-icon.png', GETDATE(), NULL),
(@knowledgeID, @productionStatusID, 'Compliance & Frameworks', 'الامتثال والأطر', 'compliance-frameworks-icon.png', GETDATE(), NULL),
(@knowledgeID, @productionStatusID, 'Metrics & Utilities', 'المقاييس والمرافق', 'metrics-utilities-icon.png', GETDATE(), NULL),
(@knowledgeID, @productionStatusID, 'Enterprise Risk Management (ERM)', '(ERM) إدارة المخاطر المؤسسية', 'erm-icon.png', GETDATE(), NULL);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE knowledge_center_resources;
-- +goose StatementEnd
