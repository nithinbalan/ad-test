-- +goose Up
-- +goose StatementBegin
INSERT INTO stakeholder_categories (title_en, title_ar, type, status) VALUES
('Department', 'قسم', 'department', 'active'),
('Vendor', 'بائع', 'vendor', 'active'),
('Client', 'عميل', 'client', 'active'),
('Individual', 'فردي', 'individual', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE stakeholder_categories;
-- +goose StatementEnd