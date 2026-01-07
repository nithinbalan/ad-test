-- +goose Up
-- +goose StatementBegin

INSERT INTO themes (
    title_en,
    title_ar,
    primary_color,
    secondary_color,
    link_color,
    sidebar_bg_color,
    is_system,
    created_by_id,
    created_at
)
VALUES
(
    'Classic Blue',
    'الأزرق الكلاسيكي',
    '#2D66E4',
    '#242D42',
    '#BC8E48',
    '#242D42',
    0,
    NULL,
    GETDATE()
),
(
    'Ocean Teal',
    'تركوازي المحيط',
    '#12687B',
    '#132F35',
    '#BC8E48',
    '#132F35',
    0,
    NULL,
    GETDATE()
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM themes
WHERE title_en IN ('Classic Blue', 'Ocean Teal');

-- +goose StatementEnd
