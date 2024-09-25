-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_finding_categories (title_en, title_ar, description_en, description_ar, type, color_code)
VALUES 
('Major Non-Compliance', 'عدم الامتثال الرئيسي', 'A significant deviation from regulations or standards.', 'انحراف كبير عن اللوائح أو المعايير.', 'major-nc', '#DC3545'),
('Minor Non-Compliance', 'عدم الامتثال البسيط', 'A minor deviation from regulations or standards.', 'انحراف بسيط عن اللوائح أو المعايير.', 'minor-nc', '#FFA500'),
('Opportunity for Improvement', 'فرصة للتحسين', 'A suggestion or recommendation for improving processes.', 'اقتراح أو توصية لتحسين العمليات.', 'ofi', '#90DD75'),
('Observation', 'ملاحظة', 'A noted issue or point of interest that does not fall under non-compliance.', 'قضية ملحوظة أو نقطة اهتمام لا تقع تحت عدم الامتثال.', 'observation', '#6C757D');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_finding_categories;
-- +goose StatementEnd
