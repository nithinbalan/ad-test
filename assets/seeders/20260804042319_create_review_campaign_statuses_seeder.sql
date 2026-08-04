-- +goose Up
-- +goose StatementBegin
INSERT INTO review_campaign_statuses (title_en, title_ar, color_code, type, status)
VALUES
    ('Planned', N'مخطط', '#6C757D', 'planned', 'active'),
    ('In Progress', N'قيد التنفيذ', '#007BFF', 'in_progress', 'active'),
    ('Submitted', N'تم التقديم', '#17A2B8', 'submitted', 'active'),
    ('Approved', N'تمت الموافقة', '#28A745', 'approved', 'active'),
    ('Closed', N'مغلق', '#343A40', 'closed', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM review_campaign_statuses
WHERE type IN ('planned', 'in_progress', 'submitted', 'approved', 'closed');
-- +goose StatementEnd
