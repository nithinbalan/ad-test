-- +goose Up
-- +goose StatementBegin
INSERT INTO kpi_achievement_logics (title_en, title_ar, type, status)
VALUES
    ('Higher is Better', N'كلما ارتفع، كان أفضل.', 'higher', 'active'),
    ('Lower is Better', N'الأقل أفضل', 'lower', 'active'),
    ('Steady or Increasing', N'ثابت أو متزايد', 'steady_or_increasing', 'active'),
    ('Increasing', N'متزايد', 'increasing', 'active'),
    ('Decreasing', N'متناقص', 'decreasing', 'active'),
    ('Steady or Decreasing', N'ثابت أو متناقص', 'steady_or_decreasing', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM kpi_achievement_logics
WHERE [type] IN (
    'higher', 'lower', 'steady_or_increasing',
    'increasing', 'decreasing', 'steady_or_decreasing'
);
-- +goose StatementEnd
