-- +goose Up
-- +goose StatementBegin
INSERT INTO user_statuses (title_en, title_ar, color_code, type, status) VALUES
('Active', 'نشيط', '#00FF00', 'active', 'active'),
('On Leave', 'في إجازة', '#FFFF00', 'onleave', 'active'),
('Inactive', 'غير نشيط', '#FF0000', 'inactive', 'active'),
('Locked', 'مقفل', '#0000FF', 'locked', 'active'),
('Retired', 'متقاعد', '#000000', 'retired', 'active'),
('Resigned', 'مستقيل', '#808080', 'resigned', 'active'),
('Deceased', 'متوفى', '#FFFFFF', 'deceased', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE user_statuses;
-- +goose StatementEnd
