-- +goose Up
-- +goose StatementBegin
INSERT INTO swot_pestel_categories (title_en, title_ar, color_code, type, status) VALUES
('Strengths', 'المقرر', '#007BFF', 'strength', 'active'),
('Weakness', 'المقرر', '#28A745', 'weakness', 'active'),
('Opportunity', 'المقرر', '#FFC107', 'opportunity', 'active'),
('Threats', 'مكتمل', '#DC3545', 'threat', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE swot_pestel_categories;
-- +goose StatementEnd