-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @RiskMatrixId INT;
SELECT @RiskMatrixId = id FROM risk_matrix WHERE type = 'erm';

DECLARE @RiskRatingIds TABLE (id INT PRIMARY KEY);
INSERT INTO @RiskRatingIds (id)
SELECT id
FROM risk_matrix_risk_ratings
WHERE risk_matrix_id = @RiskMatrixId;

DECLARE @RiskRatingValueIds TABLE (id INT PRIMARY KEY);
INSERT INTO @RiskRatingValueIds (id)
SELECT v.id
FROM risk_matrix_risk_rating_values v
INNER JOIN @RiskRatingIds rr ON rr.id = v.risk_matrix_risk_rating_id;

DECLARE @RiskRegisterIds TABLE (id INT PRIMARY KEY);
INSERT INTO @RiskRegisterIds (id)
SELECT id
FROM risk_register
WHERE risk_matrix_id = @RiskMatrixId;

DECLARE @RiskEvaluationIds TABLE (id INT PRIMARY KEY);
INSERT INTO @RiskEvaluationIds (id)
SELECT DISTINCT re.id
FROM risk_evaluations re
LEFT JOIN @RiskRatingIds rr ON rr.id = re.risk_matrix_risk_rating_id
LEFT JOIN @RiskRatingValueIds rv ON rv.id = re.risk_matrix_risk_rating_value_id
LEFT JOIN @RiskRegisterIds r ON r.id = re.risk_register_id
WHERE rr.id IS NOT NULL OR rv.id IS NOT NULL OR r.id IS NOT NULL;

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
-- 3. DELETE direct rating references
--------------------------------------------------

UPDATE controls
SET risk_matrix_risk_rating_id = NULL,
    risk_matrix_risk_rating_value_id = NULL
WHERE risk_matrix_risk_rating_id IN (SELECT id FROM @RiskRatingIds)
   OR risk_matrix_risk_rating_value_id IN (SELECT id FROM @RiskRatingValueIds);

UPDATE processes
SET risk_matrix_risk_rating_id = NULL,
    risk_matrix_risk_rating_value_id = NULL
WHERE risk_matrix_risk_rating_id IN (SELECT id FROM @RiskRatingIds)
   OR risk_matrix_risk_rating_value_id IN (SELECT id FROM @RiskRatingValueIds);

--------------------------------------------------
-- 4. DELETE ratings and values
--------------------------------------------------

DELETE FROM risk_treatment_guidelines
WHERE risk_matrix_risk_rating_id IN (SELECT id FROM @RiskRatingIds);

DELETE FROM risk_matrix_risk_rating_values
WHERE id IN (SELECT id FROM @RiskRatingValueIds);

DELETE FROM risk_matrix_risk_ratings
WHERE id IN (SELECT id FROM @RiskRatingIds);

--------------------------------------------------
-- 5. INSERT fresh ratings
--------------------------------------------------

INSERT INTO risk_matrix_risk_ratings (risk_matrix_id, title_en, title_ar, type, color_code, status)
VALUES
    (@RiskMatrixId, 'Low',       N'منخفض',      'low',       '#28A745', 'active'),
    (@RiskMatrixId, 'Medium',    N'متوسط',      'medium',    '#FFC107', 'active'),
    (@RiskMatrixId, 'High',      N'عالي',       'high',      '#FFA500', 'active'),
    (@RiskMatrixId, 'Critical',  N'حرج',        'critical',  '#FF0000', 'active');
COMMIT;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- no-op: down intentionally disabled for data seeders
SELECT 1;
-- +goose StatementEnd
