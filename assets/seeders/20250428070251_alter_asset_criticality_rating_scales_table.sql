-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @MatrixId INT;
SELECT @MatrixId = id FROM asset_criticality_matrices WHERE title_en = 'Default Asset Criticality Matrix';

-- Update Low range
UPDATE asset_criticality_matrix_rating_scales
SET score_from = 1.0, score_to = 2.0
WHERE asset_criticality_matrix_id = @MatrixId
AND type = 'low';

-- Update Medium range
UPDATE asset_criticality_matrix_rating_scales
SET score_from = 2.0, score_to = 3.5
WHERE asset_criticality_matrix_id = @MatrixId
AND type = 'medium';

-- Update High range
UPDATE asset_criticality_matrix_rating_scales
SET score_from = 3.5, score_to = 5.0
WHERE asset_criticality_matrix_id = @MatrixId
AND type = 'high';

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @MatrixId INT;
SELECT @MatrixId = id FROM asset_criticality_matrices WHERE title_en = 'Default Asset Criticality Matrix';

-- Revert Low range
UPDATE asset_criticality_matrix_rating_scales
SET score_from = 1.0, score_to = 2.0
WHERE asset_criticality_matrix_id = @MatrixId
AND type = 'low';

-- Revert Medium range
UPDATE asset_criticality_matrix_rating_scales
SET score_from = 2.1, score_to = 3.5
WHERE asset_criticality_matrix_id = @MatrixId
AND type = 'medium';

-- Revert High range
UPDATE asset_criticality_matrix_rating_scales
SET score_from = 3.6, score_to = 5.0
WHERE asset_criticality_matrix_id = @MatrixId
AND type = 'high';

COMMIT;
-- +goose StatementEnd
