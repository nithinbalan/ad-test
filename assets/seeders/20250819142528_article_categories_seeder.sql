-- +goose Up
-- +goose StatementBegin
 
INSERT INTO article_categories (title_en, title_ar, status) VALUES
('Governance', 'الحوكمة', 'active'),
('Risk', 'المخاطر', 'active'),
('Compliance', 'الامتثال', 'active');
 
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
 
DELETE FROM article_categories
WHERE title_en IN ('Governance', 'Risk', 'Compliance');
 
-- +goose StatementEnd