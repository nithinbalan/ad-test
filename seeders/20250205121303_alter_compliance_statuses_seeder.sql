-- +goose Up
-- +goose StatementBegin
UPDATE compliance_statuses 
SET score_from = 0, score_to = 50
WHERE type = 'not_implemented';

UPDATE compliance_statuses 
SET score_from = 0, score_to = 0
WHERE type = 'in_progress';

UPDATE compliance_statuses 
SET score_from = 51, score_to = 80
WHERE type = 'partially_compliant';

UPDATE compliance_statuses 
SET score_from = 81, score_to = 100
WHERE type = 'compliant';

UPDATE compliance_statuses 
SET score_from = 0, score_to = 0
WHERE type = 'exempt';

UPDATE compliance_statuses 
SET score_from = 0, score_to = 0
WHERE type = 'not_applicable';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE compliance_statuses 
SET score_from = null, score_to = null
WHERE type = 'not_implemented';

UPDATE compliance_statuses 
SET score_from = null, score_to = null
WHERE type = 'in_progress';

UPDATE compliance_statuses 
SET score_from = null, score_to = null
WHERE type = 'partially_compliant';

UPDATE compliance_statuses 
SET score_from = null, score_to = null
WHERE type = 'compliant';

UPDATE compliance_statuses 
SET score_from = null, score_to = null
WHERE type = 'exempt';

UPDATE compliance_statuses 
SET score_from = null, score_to = null
WHERE type = 'not_applicable';
-- +goose StatementEnd
