-- +goose Up
-- +goose StatementBegin

INSERT INTO knowledge_center_resource_statuses (title_en, title_ar, status)
VALUES
('New', 'جديد', 'active'),
('Coming Soon', 'قريباً', 'active'),
('Production', 'قيد التشغيل', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE knowledge_center_resource_statuses;
-- +goose StatementEnd
