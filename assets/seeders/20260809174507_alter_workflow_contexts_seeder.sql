-- +goose Up
-- +goose StatementBegin
DECLARE @identityAccessReviewModuleID INT;

SET @identityAccessReviewModuleID = (
    SELECT id
    FROM modules
    WHERE identifier = 'identity_access_review'
);

INSERT INTO workflow_contexts (
    module_id,
    title_en,
    title_ar,
    type
)
VALUES
(
    @identityAccessReviewModuleID,
    'Identity Access Review Compliance Workflow',
    N'سير عمل الامتثال لمراجعة الهوية والوصول',
    'identity_access_review_compliance'
),
(
    @identityAccessReviewModuleID,
    'Identity Access Review Cyber Compliance Workflow',
    N'سير عمل الامتثال السيبراني لمراجعة الهوية والوصول',
    'identity_access_review_cyber_compliance'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_contexts
WHERE type IN (
    'identity_access_review_compliance',
    'identity_access_review_cyber_compliance'
);
-- +goose StatementEnd
