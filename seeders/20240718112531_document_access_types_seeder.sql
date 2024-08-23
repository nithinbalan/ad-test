-- +goose Up
-- +goose StatementBegin
INSERT INTO document_access_types (title_en, title_ar, type, color_code, status) VALUES
('Public', 'عام', 'public', '#FFFFFF', 'active'),
('Shared', 'مشترك', 'shared', '#FFCC00', 'active'),
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE document_access_types;
-- +goose StatementEnd
