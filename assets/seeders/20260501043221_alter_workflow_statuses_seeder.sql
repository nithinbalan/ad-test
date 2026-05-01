-- +goose Up
-- +goose StatementBegin
INSERT INTO workflow_statuses
(title_en, title_ar, color_code, type, status)
VALUES

('Changes Requested', N'تم طلب التغييرات',
 '#FF9F45', 'changes_requested', 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_statuses
WHERE type = 'changes_requested';
-- +goose StatementEnd
