-- +goose Up
-- +goose StatementBegin
INSERT INTO control_efficiency_statuses (title_en, title_ar, color_code, type, status) VALUES
('Efficient', 'فعال', '#00FF00', 'efficient', 'active'),
('Moderate', 'معتدل', '#FFFF00', 'moderate', 'active'),
('Inefficient', 'غير فعال', '#FF0000', 'inefficient', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE control_efficiency_statuses;
-- +goose StatementEnd
