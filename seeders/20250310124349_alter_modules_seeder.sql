-- +goose Up
-- +goose StatementBegin
UPDATE modules
SET title_ar = 'مؤشرات تحليل تأثير الأعمال' , title_en = 'BIA Matrices'
WHERE identifier = 'bia_matrics';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules
SET title_ar = 'مؤشرات تحليل تأثير الأعمال' , title_en = 'BIA Matrics'
WHERE identifier = 'bia_matrics';
-- +goose StatementEnd
