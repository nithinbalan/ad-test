-- +goose Up
-- +goose StatementBegin
-- Update category names to new titles
UPDATE risk_impact_categories 
SET title_en = 'Disruption', title_ar = 'انقطاع الأعمال'
WHERE title_en = 'Business Interruption';

UPDATE risk_impact_categories 
SET title_en = 'Systems', title_ar = 'التقني'
WHERE title_en = 'Technical';

UPDATE risk_impact_categories 
SET title_en = 'Cybersecurity', title_ar = 'السيبراني'
WHERE title_en = 'Cyber';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert category names to old titles
UPDATE risk_impact_categories 
SET title_en = 'Business Interruption', title_ar = 'انقطاع الأعمال'
WHERE title_en = 'Disruption';

UPDATE risk_impact_categories 
SET title_en = 'Technical', title_ar = 'التقني'
WHERE title_en = 'Systems';

UPDATE risk_impact_categories 
SET title_en = 'Cyber', title_ar = 'السيبراني'
WHERE title_en = 'Cybersecurity';
-- +goose StatementEnd
