-- +goose Up
-- +goose StatementBegin
UPDATE swot_pestel_categories SET 
    color_code = '#1DA1F2',
    type = 'political'
WHERE title_en = 'Political';

UPDATE swot_pestel_categories SET 
    color_code = '#007BFF',
    type = 'economic'
WHERE title_en = 'Economic';

UPDATE swot_pestel_categories SET 
    color_code = '#2ECC71',
    type = 'social'
WHERE title_en = 'Social';

UPDATE swot_pestel_categories SET 
    color_code = '#F1C40F',
    type = 'technological'
WHERE title_en = 'Technological';

UPDATE swot_pestel_categories SET 
    color_code = '#E67E22',
    type = 'legal'
WHERE title_en = 'Legal';

UPDATE swot_pestel_categories SET 
    color_code = '#E74C3C',
    type = 'environmental'
WHERE title_en = 'Environmental';


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE swot_pestel_categories SET 
    color_code = '#007BFF',
    type = 'strength'
WHERE title_en = 'Political';

UPDATE swot_pestel_categories SET 
    color_code = '#28A745',
    type = 'weakness'
WHERE title_en = 'Economic';

UPDATE swot_pestel_categories SET 
    color_code = '#FFC107',
    type = 'opportunity'
WHERE title_en = 'Social';

UPDATE swot_pestel_categories SET 
    color_code = '#DC3545',
    type = 'threat'
WHERE ttitle_enype = 'Technological';

UPDATE swot_pestel_categories SET 
    color_code = '#DC3545',
    type = 'threat'
WHERE title_en = 'Legal';

UPDATE swot_pestel_categories SET 
    color_code = '#DC3545',
    type = 'threat'
WHERE title_en = 'Environmental';

-- +goose StatementEnd
