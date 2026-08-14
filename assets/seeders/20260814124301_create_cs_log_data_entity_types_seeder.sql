-- +goose Up
-- +goose StatementBegin
INSERT INTO cs_log_data_entity_types
    (title_en, title_ar, status)
VALUES
    ('Reports', N'التقارير', 'active'),
    ('Indicators', N'المؤشرات', 'active'),
    ('Threat Information', N'معلومات التهديدات', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM cs_log_data_entity_types
WHERE title_en IN ('Reports', 'Indicators', 'Threat Information');
-- +goose StatementEnd
