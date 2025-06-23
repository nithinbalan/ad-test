-- +goose Up
-- +goose StatementBegin

UPDATE modules
SET [order] = 3, title_en = 'Org. Context', title_ar = 'منظمة. سياق'
WHERE identifier = 'organizational_context';

UPDATE modules
SET [order] = 4
WHERE identifier = 'compliance';

UPDATE modules
SET [order] = 5
WHERE identifier = 'risk';

UPDATE modules
SET [order] = 6
WHERE identifier = 'audit';

UPDATE modules
SET [order] = 7
WHERE identifier = 'audit_insights';

UPDATE modules
SET [order] = 8
WHERE identifier = 'document';

UPDATE modules
SET [order] = 9
WHERE identifier = 'incident';

UPDATE modules
SET [order] = 10
WHERE identifier = 'business_continuity';

UPDATE modules
SET [order] = 11
WHERE identifier = 'management_review';

UPDATE modules
SET [order] = 12
WHERE identifier = 'asset';

UPDATE modules
SET [order] = 13
WHERE identifier = 'strategy';

UPDATE modules
SET [order] = 14
WHERE identifier = 'report';

UPDATE modules
SET [order] = 15
WHERE identifier = 'tprm';

UPDATE modules
SET [order] = 16
WHERE identifier = 'knowledge_center';

UPDATE modules
SET [order] = 17
WHERE identifier = 'setting';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules
SET [order] = 14, [title_en] = "Organizational Context"
WHERE identifier = 'organizational_context';

UPDATE modules
SET [order] = 3
WHERE identifier = 'compliance';

UPDATE modules
SET [order] = 4
WHERE identifier = 'risk';

UPDATE modules
SET [order] = 5
WHERE identifier = 'audit';

UPDATE modules
SET [order] = 6
WHERE identifier = 'sudit_insights';

UPDATE modules
SET [order] = 7
WHERE identifier = 'document';

UPDATE modules
SET [order] = 8
WHERE identifier = 'incident';

UPDATE modules
SET [order] = 9
WHERE identifier = 'business_continuity';

UPDATE modules
SET [order] = 10
WHERE identifier = 'management_review';

UPDATE modules
SET [order] = 11
WHERE identifier = 'asset';

UPDATE modules
SET [order] = 13
WHERE identifier = 'strategy';

UPDATE modules
SET [order] = 14
WHERE identifier = 'report';

UPDATE modules
SET [order] = 15
WHERE identifier = 'tprm';

UPDATE modules
SET [order] = 16
WHERE identifier = 'knowledge_center';

UPDATE modules
SET [order] = 15
WHERE identifier = 'setting';
-- +goose StatementEnd
