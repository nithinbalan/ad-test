-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

UPDATE risk_matrix_risk_ratings
SET color_code = '#28A745'
WHERE type = 'very_low';

UPDATE risk_matrix_risk_ratings
SET color_code = '#FFC107'
WHERE type = 'low';

UPDATE risk_matrix_risk_ratings
SET color_code = '#E0554D'
WHERE type = 'very_high';

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN TRANSACTION;

-- Reverting to the original color codes
UPDATE risk_matrix_risk_ratings
SET color_code = '#008000' -- Original color for very_low
WHERE type = 'very_low';

UPDATE risk_matrix_risk_ratings
SET color_code = '#FFFF00' -- Original color for low
WHERE type = 'low';

UPDATE risk_matrix_risk_ratings
SET color_code = '#B22222' -- Original color for very_high
WHERE type = 'very_high';

COMMIT;
-- +goose StatementEnd
