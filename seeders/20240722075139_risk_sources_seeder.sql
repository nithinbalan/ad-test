-- +goose Up
-- +goose StatementBegin
INSERT INTO risk_sources (title_en, title_ar, status) VALUES
('Other', 'آخر', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_sources;
-- +goose StatementEnd
