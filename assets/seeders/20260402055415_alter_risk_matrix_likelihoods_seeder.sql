-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @RiskMatrixId INT;
SELECT @RiskMatrixId = id FROM risk_matrix WHERE type = 'erm';

DECLARE @LikelihoodIds TABLE (id INT PRIMARY KEY);
INSERT INTO @LikelihoodIds (id)
SELECT id
FROM risk_matrix_likelihoods
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
LEFT JOIN @LikelihoodIds l ON l.id = re.risk_matrix_likelihood_id
LEFT JOIN @RiskRegisterIds rr ON rr.id = re.risk_register_id
WHERE l.id IS NOT NULL OR rr.id IS NOT NULL;

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

DELETE FROM control_tests
WHERE risk_register_id IN (SELECT id FROM @RiskRegisterIds);

DELETE FROM risks_strategic_intiatives
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
-- 3. DELETE likelihoods
--------------------------------------------------

DELETE FROM risk_matrix_likelihoods
WHERE id IN (SELECT id FROM @LikelihoodIds);

--------------------------------------------------
-- 4. INSERT fresh likelihoods
--------------------------------------------------

INSERT INTO risk_matrix_likelihoods
(risk_matrix_id, title_en, title_ar, score, description_en, description_ar, probability, status)
VALUES
    (@RiskMatrixId, 'Rare', N'نادر', 1,
     'Highly unlikely to occur. Expected once every 10+ years.',
     N'من غير المحتمل جدًا حدوثه. من المتوقع حدوثه مرة كل أكثر من 10 سنوات.',
     'Once every 10+ years',
     'active'),

    (@RiskMatrixId, 'Unlikely', N'غير محتمل', 2,
     'Could happen but rare. Expected once every 3-10 years.',
     N'يمكن أن يحدث ولكن نادر. من المتوقع مرة كل 3-10 سنوات.',
     'Once every 3-10 years',
     'active'),

    (@RiskMatrixId, 'Likely', N'محتمل', 3,
     'Might occur at some point. Expected once every 1-3 years.',
     N'قد يحدث في مرحلة ما. من المتوقع مرة كل 1-3 سنوات.',
     'Once every 1-3 years',
     'active'),

    (@RiskMatrixId, 'Very Likely', N'مرجح جدًا', 4,
     'Will probably occur in most circumstances. Expected to occur annually.',
     N'من المرجح حدوثه في معظم الظروف. من المتوقع حدوثه سنويًا.',
     'Annually',
     'active'),

    (@RiskMatrixId, 'Almost Certain', N'شبه مؤكد', 5,
     'Expected to occur frequently. More than once per quarter.',
     N'من المتوقع حدوثه بشكل متكرر. أكثر من مرة في الربع.',
     '> Once per quarter',
     'active');
COMMIT;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- no-op: down intentionally disabled for data seeders
SELECT 1;
-- +goose StatementEnd
