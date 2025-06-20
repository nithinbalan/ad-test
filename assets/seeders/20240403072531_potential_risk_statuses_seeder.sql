-- +goose Up
-- +goose StatementBegin
INSERT INTO dbo.potential_risk_statuses
    (title_en, title_ar, color_code, [type], [status])
VALUES
    ('New', 'جديد', '#FFC107', 'new', 'active'),
    ('Confirmed', 'مصدق عليه', '#DC3545', 'confirmed', 'active'),
    ('Dismissed', 'مرفوض', '#28A745', 'dismissed', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE potential_risk_statuses;
-- +goose StatementEnd
