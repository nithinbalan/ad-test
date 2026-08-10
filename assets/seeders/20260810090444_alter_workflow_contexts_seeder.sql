-- +goose Up
-- +goose StatementBegin
UPDATE workflow_contexts
SET
    title_en = 'Identity Access Review Cyber Security Workflow',
    title_ar = N'سير عمل الأمن السيبراني لمراجعة الهوية والوصول',
    type = 'identity_access_review_cyber_security'
WHERE type = 'identity_access_review_cyber_compliance';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE workflow_contexts
SET
    title_en = 'Identity Access Review Cyber Compliance Workflow',
    title_ar = N'سير عمل الامتثال السيبراني لمراجعة الهوية والوصول',
    type = 'identity_access_review_cyber_compliance'
WHERE type = 'identity_access_review_cyber_security';
-- +goose StatementEnd
