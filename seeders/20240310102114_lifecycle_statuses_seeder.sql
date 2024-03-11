-- +goose Up
-- +goose StatementBegin
INSERT INTO lifecycle_statuses (title_en, title_ar, color_code, type, status) VALUES
('Active', 'نشط', '#00FF00', 'active', 'active'),
('Expired', 'منتهي الصلاحية', '#0000FF', 'expired', 'active'),
('Archived', 'أرشفة', '#808080', 'archived', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE lifecycle_statuses;
-- +goose StatementEnd
