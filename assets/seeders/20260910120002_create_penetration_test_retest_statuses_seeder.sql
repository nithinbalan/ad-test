-- +goose Up
-- +goose StatementBegin
INSERT INTO penetration_test_retest_statuses (title_en, title_ar, color_code, type, status)
VALUES ('Draft', N'مسودة', '#9CA3AF', 'draft', 'active'), ('In Progress', N'قيد التنفيذ', '#F59E0B', 'in-progress', 'active'), ('Completed', N'مكتمل', '#10B981', 'completed', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE penetration_test_retest_statuses;
-- +goose StatementEnd
