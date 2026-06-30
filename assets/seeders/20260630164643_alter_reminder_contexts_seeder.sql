-- +goose Up
-- +goose StatementBegin
DECLARE @ComplianceModuleID INT;
DECLARE @DocumentModuleID INT;
DECLARE @IncidentModuleID INT;
DECLARE @MeetingModuleID INT;
DECLARE @ProjectModuleID INT;
DECLARE @RiskModuleID INT;
DECLARE @AssetModuleID INT;
DECLARE @StakeholderModuleID INT;

SELECT @ComplianceModuleID = id FROM modules WHERE identifier = 'compliance';
SELECT @DocumentModuleID = id FROM modules WHERE identifier = 'document';
SELECT @IncidentModuleID = id FROM modules WHERE identifier = 'incident';
SELECT @MeetingModuleID = id FROM modules WHERE identifier = 'meeting';
SELECT @ProjectModuleID = id FROM modules WHERE identifier = 'organization';
SELECT @RiskModuleID = id FROM modules WHERE identifier = 'risk';
SELECT @AssetModuleID = id FROM modules WHERE identifier = 'asset';
SELECT @StakeholderModuleID = id FROM modules WHERE identifier = 'organizational_context';

INSERT INTO reminder_contexts (type, title_en, title_ar, status, module_id)
VALUES
(
    'threat-action-plan-reminder',
    'Threat Action Plan Reminder',
    N'تذكير خطة عمل التهديد',
    'active',
    @AssetModuleID
),
(
    'compliance-deadline-due',
    'Compliance Deadline Due',
    N'استحقاق موعد الامتثال',
    'active',
    @ComplianceModuleID
),
(
    'control-action-plan-reminder',
    'Control Action Plan Reminder',
    N'تذكير خطة عمل الضبط',
    'active',
    @ComplianceModuleID
),
(
    'sla-action-plan-reminder',
    'SLA Action Plan Reminder',
    N'تذكير خطة عمل اتفاقية مستوى الخدمة',
    'active',
    @ComplianceModuleID
),
(
    'document-review-reminder',
    'Document Review Reminder',
    N'تذكير مراجعة الوثيقة',
    'active',
    @DocumentModuleID
),
(
    'incident-action-plan-reminder',
    'Incident Action Plan Reminder',
    N'تذكير خطة عمل الحادث',
    'active',
    @IncidentModuleID
),
(
    'meeting-action-plan-reminder',
    'Meeting Action Plan Reminder',
    N'تذكير خطة عمل الاجتماع',
    'active',
    @MeetingModuleID
),
(
    'stakeholder-action-plan-reminder',
    'Stakeholder Action Plan Reminder',
    N'تذكير خطة عمل أصحاب المصلحة',
    'active',
    @StakeholderModuleID
),
(
    'project-action-plan-reminder',
    'Project Action Plan Reminder',
    N'تذكير خطة عمل المشروع',
    'active',
    @ProjectModuleID
),
(
    'risk-action-plan-due',
    'Risk Action Plan Due',
    N'استحقاق خطة عمل المخاطر',
    'active',
    @RiskModuleID
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM reminder_contexts
WHERE type IN (
    'threat-action-plan-reminder',
    'compliance-deadline-due',
    'control-action-plan-reminder',
    'sla-action-plan-reminder',
    'document-review-reminder',
    'incident-action-plan-reminder',
    'meeting-action-plan-reminder',
    'stakeholder-action-plan-reminder',
    'project-action-plan-reminder',
    'risk-action-plan-due'
);
-- +goose StatementEnd