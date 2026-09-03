-- +goose Up
-- +goose StatementBegin
INSERT INTO planned_advisory_sources (title_en, title_ar, status) VALUES
    ('NCA', N'الهيئة الوطنية للأمن السيبراني', 'active'),
    ('National CERT', N'المركز الوطني للاستجابة للحوادث الإلكترونية', 'active'),
    ('ministry', N'وزارة', 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM planned_advisory_sources
WHERE title_en IN ('NCA', 'National CERT', 'ministry');
-- +goose StatementEnd
