-- +goose Up
-- +goose StatementBegin
UPDATE compliance_statuses 
SET score_from = 0, 
    score_to = 24,
    base_score = 0,
    title_en = 'Non Compliant',
    title_ar = 'غير متوافق'
WHERE type = 'not_implemented';

UPDATE compliance_statuses 
SET score_from = 25, score_to = 75,
    base_score = 25
WHERE type = 'partially_compliant';

UPDATE compliance_statuses 
SET score_from = 100, score_to = 100,
    base_score = 100
WHERE type = 'compliant';

UPDATE compliance_statuses 
SET score_from = 0, score_to = 0,
    base_score = 0
WHERE type = 'not_applicable';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE compliance_statuses 
SET score_from = 0, 
    score_to = 50,
    base_score = 0,
    title_en = 'Not Implemented',
    title_ar = 'لم تنفذ'
WHERE type = 'not_implemented';

UPDATE compliance_statuses 
SET score_from = 50, score_to = 80,
    base_score = 50
WHERE type = 'partially_compliant';

UPDATE compliance_statuses 
SET score_from = 80, score_to = 100,
    base_score = 100
WHERE type = 'compliant';

UPDATE compliance_statuses 
SET score_from = 0, score_to = 0
WHERE type = 'not_applicable';
-- +goose StatementEnd
