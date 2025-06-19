-- +goose Up
-- +goose StatementBegin
INSERT INTO third_party_question_types (title_en, title_ar, type, status) VALUES
('Yes/No', 'نعم / لا', 'yes-no', 'active'),
('Dropdown (Single)', 'قائمة منسدلة (مفردة)', 'dropdown-single', 'active'),
('Dropdown (Multiple)', 'قائمة منسدلة (متعددة)', 'dropdown-multiple', 'active'),
('Rating Scale', 'مقياس التقييم', 'rating-scale', 'active'),
('Single Text', 'نص مفرد', 'single-text', 'active'),
('Text Area', 'منطقة نصية', 'text-area', 'active'),
('File Upload', 'تحميل ملف', 'file-upload', 'active'),
('Date Picker', 'منتقي التاريخ', 'date-picker', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM third_party_question_types;
-- +goose StatementEnd
