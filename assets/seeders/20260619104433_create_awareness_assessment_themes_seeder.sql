-- +goose Up
-- +goose StatementBegin
INSERT INTO awareness_assessment_themes (title_en, title_ar)
VALUES
('Phishing', N'التصيد الاحتيالي'),
('Passwords', N'كلمات المرور'),
('Data Protection', N'حماية البيانات'),
('Clarity', N'الوضوح'),
('Relevance', N'الملاءمة'),
('Accessibility', N'سهولة الوصول'),
('Infographics', N'الرسوم البيانية'),
('Applicability', N'قابلية التطبيق');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM awareness_assessment_themes
-- +goose StatementEnd
