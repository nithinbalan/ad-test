-- +goose Up
-- +goose StatementBegin
INSERT INTO asset_statuses (title_en, title_ar, color_code, type)
VALUES
    ('Discovered', N'تم الاكتشاف', '#007BFF', 'discovered'),
    ('Unavailable', N'غير متاح', '#6C757D', 'unavailable'),
    ('Decommissioned', N'خارج الخدمة', '#FEC107', 'decommissioned'),
    ('Missing', N'مفقود', '#FD7E14', 'missing'),
    ('Under Review', N'قيد المراجعة', '#90DD75', 'under-review');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM asset_statuses
WHERE type IN (
    'discovered',
    'unavailable',
    'decommissioned',
    'missing',
    'under-review'
);
-- +goose StatementEnd