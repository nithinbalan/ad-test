-- +goose Up
-- +goose StatementBegin
INSERT INTO control_efficiency_statuses (title_en, title_ar, color_code, type, status, score_from, score_to) VALUES
('Efficient', 'فعال', '#00FF00', 'efficient', 'active', 80, 100),
('Moderate', 'معتدل', '#FFFF00', 'moderate', 'active', 31, 79),
('Inefficient', 'غير فعال', '#FF0000', 'inefficient', 'active', 0, 30);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE control_efficiency_statuses;
-- +goose StatementEnd
