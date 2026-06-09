-- +goose Up
-- +goose StatementBegin
-- Update Critical: 4.0 – 5.0, Red
UPDATE compliance_assessment_ratings
SET score_from = 4.0,
    score_to = 5.0,
    color_code = '#B22222'
WHERE type = 'critical';

-- Update High: 3.0 – 3.99, Orange
UPDATE compliance_assessment_ratings
SET score_from = 3.0,
    score_to = 3.99,
    color_code = '#FD7E14'
WHERE type = 'high';

-- Update Medium: 2.0 – 2.99, Amber
UPDATE compliance_assessment_ratings
SET score_from = 2.0,
    score_to = 2.99,
    color_code = '#FFC107'
WHERE type = 'medium';

-- Update Low: 1.01 – 1.99, Green
UPDATE compliance_assessment_ratings
SET score_from = 1.01,
    score_to = 1.99,
    color_code = '#008000'
WHERE type = 'low';

-- Update Closed: Residual = exactly 1.0, Grey
UPDATE compliance_assessment_ratings
SET score_from = 1.0,
    score_to = 1.0,
    color_code = '#6C757D'
WHERE type = 'closed';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert to original values
UPDATE compliance_assessment_ratings
SET score_from = 3.00,
    score_to = 4.00,
    color_code = '#B22222'
WHERE type = 'critical';

UPDATE compliance_assessment_ratings
SET score_from = 2.00,
    score_to = 2.99,
    color_code = '#FF0000'
WHERE type = 'high';

UPDATE compliance_assessment_ratings
SET score_from = 1.00,
    score_to = 1.99,
    color_code = '#FFA500'
WHERE type = 'medium';

UPDATE compliance_assessment_ratings
SET score_from = 0.01,
    score_to = 0.99,
    color_code = '#90DD75'
WHERE type = 'low';

UPDATE compliance_assessment_ratings
SET score_from = 0.00,
    score_to = 0.00,
    color_code = '#008000'
WHERE type = 'closed';
-- +goose StatementEnd
