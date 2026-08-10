-- +goose Up
-- +goose StatementBegin
DECLARE @identityAccessReviewCyberSecurityContextID INT;

SET @identityAccessReviewCyberSecurityContextID = (
    SELECT id
    FROM workflow_contexts
    WHERE type = 'identity_access_review_cyber_security'
);

UPDATE workflows
SET
    type = 'identity-access-review-cyber-security-workflow',
    title_en = 'Identity Access Review Cyber Security Workflow',
    title_ar = N'سير عمل الأمن السيبراني لمراجعة الهوية والوصول'
WHERE workflow_context_id = @identityAccessReviewCyberSecurityContextID
    AND type = 'identity-access-review-cyber-compliance-workflow';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @identityAccessReviewCyberSecurityContextID INT;

SET @identityAccessReviewCyberSecurityContextID = (
    SELECT id
    FROM workflow_contexts
    WHERE type = 'identity_access_review_cyber_security'
);

UPDATE workflows
SET
    type = 'identity-access-review-cyber-compliance-workflow',
    title_en = 'Identity Access Review Cyber Compliance Workflow',
    title_ar = N'سير عمل الامتثال السيبراني لمراجعة الهوية والوصول'
WHERE workflow_context_id = @identityAccessReviewCyberSecurityContextID
    AND type = 'identity-access-review-cyber-security-workflow';
-- +goose StatementEnd
