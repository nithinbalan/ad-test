-- +goose Up
-- +goose StatementBegin
INSERT INTO account_decisions (
    title_en,
    title_ar,
    color_code,
    type,
    status
)
VALUES
('Keep', N'إبقاء', '#28A745', 'keep', 'active'),
('Revoke', N'إلغاء', '#DC3545', 'revoke', 'active'),
('Modify', N'تعديل', '#FFC107', 'modify', 'active'),
('Investigate', N'تحقيق', '#17A2B8', 'investigate', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM account_decisions
WHERE type IN ('keep', 'revoke', 'modify', 'investigate');
-- +goose StatementEnd
