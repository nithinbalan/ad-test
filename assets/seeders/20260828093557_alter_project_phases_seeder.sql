-- +goose Up
-- +goose StatementBegin
INSERT INTO project_phases (title_en, title_ar, type, status) VALUES
('Initiation', N'البدء', 'initiation', 'active'),
('Planning', N'التخطيط', 'planning', 'active'),
('Execution', N'التنفيذ', 'execution', 'active'),
('Go-Live', N'الإطلاق', 'go_live', 'active'),
('Closure', N'الإغلاق', 'closure', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM project_phases
WHERE type IN ('initiation', 'planning', 'execution', 'go_live', 'closure');
-- +goose StatementEnd