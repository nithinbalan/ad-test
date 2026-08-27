-- +goose Up
-- +goose StatementBegin
INSERT INTO penetration_test_execution_statuses (title_en, title_ar, type, status)
VALUES ('Planned', N'مخطط', 'planned', 'active'), ('Environment Preparation', N'إعداد البيئة', 'environment-preparation', 'active'), ('In Progress', N'قيد التنفيذ', 'in-progress', 'active'), ('Completed', N'مكتمل', 'completed', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE penetration_test_execution_statuses;
-- +goose StatementEnd
