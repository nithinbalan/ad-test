-- +goose Up
-- +goose StatementBegin
UPDATE swot_pestel_categories
SET is_swot = 1, is_pestel = 0
WHERE id IN (1, 2, 3, 4);

INSERT INTO swot_pestel_categories (title_en, title_ar, color_code, type, is_swot, is_pestel, status) VALUES
('Political', 'المقرر', '#007BFF', 'strength', 0, 1, 'active'),
('Economic', 'المقرر', '#28A745', 'weakness', 0, 1, 'active'),
('Social', 'المقرر', '#FFC107', 'opportunity', 0, 1, 'active'),
('Technological', 'مكتمل', '#DC3545', 'threat', 0, 1, 'active'),
('Legal', 'مكتمل', '#DC3545', 'threat', 0, 1, 'active'),
('Environmental', 'مكتمل', '#DC3545', 'threat', 0, 1, 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE swot_pestel_categories;
-- +goose StatementEnd
