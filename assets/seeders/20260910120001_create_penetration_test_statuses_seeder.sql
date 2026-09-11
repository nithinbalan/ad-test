-- +goose Up
-- +goose StatementBegin
INSERT INTO penetration_test_statuses (title_en, title_ar, color_code, type, status)
VALUES ('Draft', N'مسودة', '#9CA3AF', 'draft', 'active'), ('Under Retest', N'قيد إعادة الاختبار', '#F59E0B', 'under retest', 'active'), ('Closed', N'مغلق', '#10B981', 'closed', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE penetration_test_statuses;
-- +goose StatementEnd
