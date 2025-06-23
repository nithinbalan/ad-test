-- +goose Up
-- +goose StatementBegin
INSERT INTO control_implementation_statuses (title_en, title_ar, color_code, type, status) VALUES
('Implemented', 'منفذ', '#28A745', 'implemented', 'active'),
('Not Implemented', 'غير منفذ', '#DC3545', 'not-implemented', 'active'),
('Partially Implemented', 'منفذ جزئيا', '#FFC107', 'partially-implemented', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE control_implementation_statuses;
-- +goose StatementEnd
