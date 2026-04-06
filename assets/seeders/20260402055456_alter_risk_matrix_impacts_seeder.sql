-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @RiskMatrixId INT;
SELECT @RiskMatrixId = id FROM risk_matrix WHERE type = 'erm';

DECLARE @ImpactIds TABLE (id INT PRIMARY KEY);
INSERT INTO @ImpactIds (id)
SELECT id
FROM risk_matrix_impacts
WHERE risk_matrix_id = @RiskMatrixId;

DECLARE @RiskRegisterIds TABLE (id INT PRIMARY KEY);
INSERT INTO @RiskRegisterIds (id)
SELECT id
FROM risk_register
WHERE risk_matrix_id = @RiskMatrixId;

DECLARE @RiskEvaluationIds TABLE (id INT PRIMARY KEY);
INSERT INTO @RiskEvaluationIds (id)
SELECT DISTINCT re.id
FROM risk_evaluations re
LEFT JOIN @ImpactIds i ON i.id = re.risk_matrix_impact_id
LEFT JOIN @RiskRegisterIds rr ON rr.id = re.risk_register_id
WHERE i.id IS NOT NULL OR rr.id IS NOT NULL;

--------------------------------------------------
-- 1. DELETE risk_evaluations and children
--------------------------------------------------

DELETE FROM risk_impact_analyses
WHERE risk_evaluation_id IN (SELECT id FROM @RiskEvaluationIds);

DELETE FROM risk_appetite_alignments
WHERE risk_evaluation_id IN (SELECT id FROM @RiskEvaluationIds);

DELETE FROM risk_evaluations
WHERE id IN (SELECT id FROM @RiskEvaluationIds);

--------------------------------------------------
-- 2. DELETE risk_register dependents
--------------------------------------------------

DELETE FROM risk_workflow_reviews
WHERE risk_workflow_id IN (
    SELECT id FROM risk_workflows
    WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds)
);

DELETE FROM risk_workflow_participants
WHERE risk_workflow_id IN (
    SELECT id FROM risk_workflows
    WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds)
);

DELETE FROM violation_report_risks
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_workflows
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM control_tests
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_documents
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_root_cause_analysis
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_action_plan_status_update_documents
WHERE risk_action_plan_status_update_id IN (
    SELECT id FROM risk_action_plan_status_updates
    WHERE risk_action_plan_id IN (
        SELECT id FROM risk_action_plans
        WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds)
    )
);

DELETE FROM risk_action_plan_status_updates
WHERE risk_action_plan_id IN (
    SELECT id FROM risk_action_plans
    WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds)
);

DELETE FROM risk_action_plan_comment_documents
WHERE risk_action_plan_comment_id IN (
    SELECT id FROM risk_action_plan_comments
    WHERE risk_action_plan_id IN (
        SELECT id FROM risk_action_plans
        WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds)
    )
);

DELETE FROM risk_action_plan_comments
WHERE risk_action_plan_id IN (
    SELECT id FROM risk_action_plans
    WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds)
);

DELETE FROM risk_action_plan_controls
WHERE risk_action_plan_id IN (
    SELECT id FROM risk_action_plans
    WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds)
);

DELETE FROM risk_action_plans
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_control_analyses
WHERE risk_control_id IN (
    SELECT id FROM risk_controls
    WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds)
);

DELETE FROM risk_controls
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_compliance_requirements
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM incident_risks
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_processes
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_strategic_objectives
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_slas
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_kris
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

UPDATE auditable_items
SET risk_register_id = NULL
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_timelines
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_associated_documents
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM meeting_risks
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM audit_plan_auditable_item_risk_controls
WHERE audit_plan_auditable_item_risk_id IN (
    SELECT id FROM audit_plan_auditable_item_risks
    WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds)
);

DELETE FROM audit_plan_auditable_item_risks
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_threats
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_assets
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_comments
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM kpi_risks
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM project_risks
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM third_party_risks
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM audit_finding_risks
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM bcp_risks
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_closures
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM policy_exception_risks
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risk_register
WHERE id IN (SELECT id FROM @RiskRegisterIds);

--------------------------------------------------
-- 3. DELETE impacts
--------------------------------------------------

DELETE FROM risk_matrix_impact_guidelines
WHERE risk_matrix_impact_id IN (SELECT id FROM @ImpactIds);

DELETE FROM risk_matrix_impacts
WHERE id IN (SELECT id FROM @ImpactIds);

--------------------------------------------------
-- 4. INSERT fresh impacts
--------------------------------------------------

INSERT INTO risk_matrix_impacts (risk_matrix_id, title_en, title_ar, score, description_en, description_ar, status)
VALUES
    (@RiskMatrixId, 'Very Low', N'منخفض جدًا', 1,
     'Limited or isolated impact. Minor data leak, minor audit finding, few customer complaints, minimal operational disruption without customer impact.',
     N'تأثير محدود أو معزول. تسرب بيانات محدود، نتيجة تدقيق بسيطة، قليل من شكاوى العملاء، اضطراب تشغيلي طفيف دون تأثير على العملاء.',
     'active'),
    (@RiskMatrixId, 'Low', N'منخفض', 2,
     'Limited or isolated impact. Minor data leak, minor audit finding, few customer complaints, minimal operational disruption without customer impact.',
     N'تأثير محدود أو معزول. تسرب بيانات محدود، نتيجة تدقيق بسيطة، قليل من شكاوى العملاء، اضطراب تشغيلي طفيف دون تأثير على العملاء.',
     'active'),
    (@RiskMatrixId, 'Moderate', N'متوسط', 3,
     'Moderate impact. Exposure of internal/confidential data, reportable compliance incident, sector-wide reputational concern, temporary or localized service interruption.',
     N'تأثير متوسط. كشف بيانات داخلية أو سرية، حادثة امتثال قابلة للإبلاغ، قلق سمعي على مستوى القطاع، انقطاع خدمة مؤقت أو محلي.',
     'active'),
    (@RiskMatrixId, 'High', N'عالي', 4,
     'High impact. Partial exposure of sensitive data impacting critical systems, formal regulatory warnings, national press or social media backlash, extended service outage under 24 hours.',
     N'تأثير عالٍ. كشف جزئي لبيانات حساسة تؤثر على الأنظمة الحيوية، تحذيرات رقابية رسمية، تداول صحفي وطني أو تأثير سلبي على وسائل التواصل، انقطاع ممتد في الخدمة أقل من 24 ساعة.',
     'active'),
    (@RiskMatrixId, 'Very High', N'عالية جداً', 5,
     'Very high impact. Full exposure of highly sensitive customer or employee data, regulatory penalties or license implications, severe media coverage causing loss of stakeholder trust, bank-wide unavailability exceeding 24 hours or loss of critical functions.',
     N'تأثير عالٍ جداً. كشف كامل لبيانات العملاء أو الموظفين الحساسة للغاية، عقوبات تنظيمية أو آثار على الترخيص، تغطية إعلامية حادة تؤدي إلى فقدان ثقة أصحاب المصلحة، عدم توفر البنك بالكامل لأكثر من 24 ساعة أو فقدان الوظائف الحيوية.',
     'active');
COMMIT;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- no-op: down intentionally disabled for data seeders
SELECT 1;
-- +goose StatementEnd
