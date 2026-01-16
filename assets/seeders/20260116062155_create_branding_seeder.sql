-- +goose Up
-- +goose StatementBegin
INSERT INTO branding (
    login_page_title_en,
    login_page_title_ar,
    login_page_subtitle_en,
    login_page_subtitle_ar,
    login_page_bg_image_path,
    favicon_image_path,
    loader_gif_path,
    org_logo_image_path,
    created_by_id,
    created_at
)
VALUES (
    'Governance, Risk and Compliance (GRC)',
    '(جي آر سي) الحوكمة والمخاطر والامتثال',
    'Welcome to the most powerful, modern AI-integrated GRC (Governance, Risk, & Compliance) platform.',
    'مرحبًا بك في أقوى وأحدث منصة للحوكمة والمخاطر والامتثال (الحوكمة والمخاطر والامتثال) المتكاملة مع الذكاء الاصطناعي.',
    NULL,
    NULL,
    NULL,
    NULL,
    1,
    GETDATE()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE branding;
-- +goose StatementEnd