-- +goose Up
-- +goose StatementBegin
INSERT INTO project_classifications (title_en, title_ar, status) VALUES
('Strategic', N'استراتيجي', 'active'),
('Operational', N'تشغيلي', 'active'),
('Compliance', N'امتثال', 'active'),
('Enhancement', N'تحسيني', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM project_classifications
WHERE title_en IN ('Strategic', 'Operational', 'Compliance', 'Enhancement');
-- +goose StatementEnd
