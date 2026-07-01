-- +goose Up
-- +goose StatementBegin
INSERT INTO asset_environments (
    title_en,
    title_ar,
    type
)
VALUES
    ('Production', N'الإنتاج', 'production'),
    ('Test', N'الاختبار', 'test'),
    ('Development', N'التطوير', 'development');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM asset_environments;
-- +goose StatementEnd