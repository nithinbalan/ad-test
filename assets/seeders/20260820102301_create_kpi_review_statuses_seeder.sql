-- +goose Up
-- +goose StatementBegin
INSERT INTO kpi_review_statuses (
    type,
    title_en,
    title_ar,
    color_code,
    status
)
VALUES
    ('updated', 'Updated', N'تم التحديث', '#00AF50', 'active'),
    ('awaiting-update', 'Pending', N'قيد الانتظار', '#FD7E14', 'active'),
    ('overdue', 'Overdue', N'متأخر', '#DC3545', 'active'),
    ('submitted', 'Submitted', N'مُقَدَّم', '#17A2B8', 'active'),
    ('approved', 'Approved', N'موافقة', '#28A745', 'active'),
    ('reverted', 'Reverted', N'عادت', '#FF9F45', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM kpi_review_statuses
WHERE type IN (
    'updated',
    'awaiting-update',
    'overdue',
    'submitted',
    'approved',
    'reverted'
);
-- +goose StatementEnd
