-- +goose Up
-- +goose StatementBegin
INSERT INTO themes (
    title_en,
    title_ar,
    description_en,
    description_ar,
    primary_color,
    secondary_color,
    link_color,
    sidebar_bg_color,
    is_system,
    created_by_id,
    created_at
)
VALUES (
    'Jethur',
    'جيثور',
    'Primary Jethur theme using brand-aligned green and gold tones, reflecting governance, stability, and trust.',
    'السمة الأساسية لجيثور باستخدام درجات اللون الأخضر والذهبي المتوافقة مع العلامة التجارية، مما يعكس الحوكمة والاستقرار والثقة.',
    '#165D31',
    '#BC8E48',
    '#B6630A',
    '#222222',
    1,
    NULL,         
    GETDATE()
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM themes
WHERE title_en = 'Jethur';
-- +goose StatementEnd
