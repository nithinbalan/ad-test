-- +goose Up
-- +goose StatementBegin
INSERT INTO communication_action_plan_statuses (title_en, title_ar, color_code, type, status) VALUES
('Sent', 'المقرر', '#90DD75', 'sent', 'active'),
('Pending', 'مكتمل', '#6C757D', 'pending', 'active');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE communication_action_plan_statuses;
-- +goose StatementEnd