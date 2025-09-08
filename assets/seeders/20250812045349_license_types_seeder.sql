-- +goose Up
-- +goose StatementBegin
INSERT INTO license_types (title_en, title_ar,color_code,type, status) VALUES
('Active', 'نشيط', '#28A745', 'active', 'active'),
('Passive', 'سلبي', '#DC3545', 'passive', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
