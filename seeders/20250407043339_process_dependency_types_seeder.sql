-- +goose Up
-- +goose StatementBegin
INSERT INTO process_dependency_types (title_en, title_ar, type, status) VALUES
('Internal', 'داخلي', 'internal', 'active'),
('External', 'خارجي', 'external', 'active');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE process_dependency_types;
-- +goose StatementEnd