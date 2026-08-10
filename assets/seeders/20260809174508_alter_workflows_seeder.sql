-- +goose Up
-- +goose StatementBegin
DECLARE @identityAccessReviewComplianceContextID INT;
DECLARE @identityAccessReviewCyberComplianceContextID INT;

SET @identityAccessReviewComplianceContextID = (
    SELECT id
    FROM workflow_contexts
    WHERE [type] = 'identity_access_review_compliance'
);

SET @identityAccessReviewCyberComplianceContextID = (
    SELECT id
    FROM workflow_contexts
    WHERE [type] = 'identity_access_review_cyber_compliance'
);

INSERT INTO workflows (
    type,
    status,
    created_by_id,
    created_at,
    title_en,
    title_ar,
    workflow_context_id
)
VALUES
(
    'identity-access-review-compliance-workflow',
    'active',
    1,
    GETDATE(),
    'Identity Access Review Compliance Workflow',
    N'سير عمل الامتثال لمراجعة الهوية والوصول',
    @identityAccessReviewComplianceContextID
),
(
    'identity-access-review-cyber-compliance-workflow',
    'active',
    1,
    GETDATE(),
    'Identity Access Review Cyber Compliance Workflow',
    N'سير عمل الامتثال السيبراني لمراجعة الهوية والوصول',
    @identityAccessReviewCyberComplianceContextID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflows
WHERE type IN (
    'identity-access-review-compliance-workflow',
    'identity-access-review-cyber-compliance-workflow'
);
-- +goose StatementEnd
