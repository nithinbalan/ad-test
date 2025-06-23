-- +goose Up
-- +goose StatementBegin
INSERT INTO asset_statuses (title_en, title_ar, color_code, type, status) VALUES
('Active', 'نشط', '#90DD75', 'active', 'active'),
('Available', 'متاح', '#17A2B8', 'available', 'active'),
('Maintenance', 'صيانة', '#FFC107', 'maintenance', 'active'),
('Disposed', 'تالف', '#E0554D', 'disposed', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE asset_statuses;
-- +goose StatementEnd
