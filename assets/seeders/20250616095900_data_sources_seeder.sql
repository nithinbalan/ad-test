-- +goose Up
-- +goose StatementBegin
INSERT INTO data_sources (title_en, title_ar,type, status)
VALUES
    ('User Input', 'إدخال المستخدم', 'user_input', 'active'),
    ('Internal Module', 'الوحدة الداخلية', 'internal_module', 'active'),
    ('External System(API)', 'النظام الخارجي (API)', 'external_api', 'active'),
    ('AI-Generated', 'تم إنشاؤه بواسطة الذكاء الاصطناعي', 'ai_generated', 'active'),
    ('Uploaded File', 'الملف الذي تم تحميله', 'file_upload', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE data_sources;
-- +goose StatementEnd
