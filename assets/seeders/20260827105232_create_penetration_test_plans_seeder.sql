-- +goose Up
-- +goose StatementBegin
INSERT INTO penetration_test_plans (title_en, title_ar, status)
VALUES ('Pentest', N'اختبار الاختراق', 'active'), ('Vulnerability management', N'إدارة الثغرات', 'active'), ('Access review plan', N'خطة مراجعة الوصول', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE penetration_test_plans;
-- +goose StatementEnd
