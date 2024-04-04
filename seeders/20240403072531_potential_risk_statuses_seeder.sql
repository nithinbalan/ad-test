-- +goose Up
-- +goose StatementBegin
INSERT INTO dbo.potential_risk_statuses (title_en, title_ar, color_code, [type], [status]) VALUES
('New', 'جديد', '#c13ce6', 'new', 'active'),
('Confirmed', 'مصدق عليه', '#e32228', 'confirmed', 'active'),
('Dismissed', 'مرفوض', '#008000', 'dismissed', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE potential_risk_statuses;
-- +goose StatementEnd
