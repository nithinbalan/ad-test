-- +goose Up
-- +goose StatementBegin
INSERT INTO project_impact_assessment_types (
    title_en, title_ar, type, status
) VALUES
('Timeline', N'الجدول الزمني', 'timeline', 'active'),
('Budget', N'الميزانية', 'budget', 'active'),
('Resources', N'الموارد', 'resources', 'active'),
('Technical Complexity', N'التعقيد الفني', 'technical_complexity', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM project_impact_assessment_types
WHERE type IN ('timeline', 'budget', 'resources', 'technical_complexity');
-- +goose StatementEnd
