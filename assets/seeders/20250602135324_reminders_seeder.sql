-- +goose Up
-- +goose StatementBegin
DECLARE @organizationModuleID INT;
DECLARE @assetModuleID INT;
DECLARE @complianceModuleID INT;
DECLARE @auditModuleID INT;
DECLARE @documentModuleID INT;
DECLARE @auditInsightModuleID INT;
DECLARE @riskModuleID INT;
DECLARE @incidentModuleID INT;
DECLARE @businessContinuityModuleID INT;
DECLARE @managementReviewModuleID INT;
DECLARE @strategyModuleID INT;
DECLARE @organizationalContextModuleID INT;
DECLARE @tprmModuleID INT;
DECLARE @reportModuleID INT;
DECLARE @knowledgeCenterModuleID INT;

SET @organizationModuleID = (SELECT id FROM modules WHERE identifier = 'organization');
SET @assetModuleID = (SELECT id FROM modules WHERE identifier = 'asset');
SET @complianceModuleID = (SELECT id FROM modules WHERE identifier = 'compliance');
SET @auditModuleID = (SELECT id FROM modules WHERE identifier = 'audit');
SET @documentModuleID = (SELECT id FROM modules WHERE identifier = 'document');
SET @auditInsightModuleID = (SELECT id FROM modules WHERE identifier = 'audit_insights');
SET @riskModuleID = (SELECT id FROM modules WHERE identifier = 'risk');
SET @incidentModuleID = (SELECT id FROM modules WHERE identifier = 'incident');
SET @businessContinuityModuleID = (SELECT id FROM modules WHERE identifier = 'business_continuity');
SET @managementReviewModuleID = (SELECT id FROM modules WHERE identifier = 'management_review');
SET @strategyModuleID = (SELECT id FROM modules WHERE identifier = 'strategy');
SET @organizationalContextModuleID = (SELECT id FROM modules WHERE identifier = 'organizational_context');
SET @tprmModuleID = (SELECT id FROM modules WHERE identifier = 'tprm');
SET @reportModuleID = (SELECT id FROM modules WHERE identifier = 'report');
SET @knowledgeCenterModuleID = (SELECT id FROM modules WHERE identifier = 'knowledge_center');


INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@organizationModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@assetModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@complianceModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@auditModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@documentModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@auditInsightModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@riskModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@incidentModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@businessContinuityModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@managementReviewModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@strategyModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@organizationalContextModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@tprmModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@reportModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

INSERT INTO reminders (module_id,title_en, title_ar,type, status) VALUES
(@knowledgeCenterModuleID,'User Tasks', 'مهام المستخدم','user-task', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE reminders;
-- +goose StatementEnd