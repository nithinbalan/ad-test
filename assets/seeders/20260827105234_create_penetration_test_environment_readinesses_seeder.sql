-- +goose Up
-- +goose StatementBegin
INSERT INTO penetration_test_environment_readinesses (title_en, title_ar, status)
VALUES ('Delivered', N'تم التسليم', 'active'), ('In Progress', N'قيد التنفيذ', 'active'), ('No Response', N'لا توجد استجابة', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE penetration_test_environment_readinesses;
-- +goose StatementEnd
