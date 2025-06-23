-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @MatrixId INT;
SELECT @MatrixId = id FROM asset_criticality_matrices WHERE title_en = 'Default Asset Criticality Matrix';

-- Insert common ratings for the matrix with specified color codes
INSERT INTO asset_criticality_matrix_factor_ratings (asset_criticality_matrix_id, title_en, title_ar, score, type, color_code)
VALUES
    (@MatrixId, 'Very Low', N'منخفض جداً', 1, 'very_low', '#90DD75'), -- Light Green
    (@MatrixId, 'Low', N'منخفض', 2, 'low', '#28A745'), -- Green
    (@MatrixId, 'Medium', N'متوسط', 3, 'medium', '#FFC107'), -- Yellow
    (@MatrixId, 'High', N'عالي', 4, 'high', '#FFA500'), -- Dark Yellow
    (@MatrixId, 'Very High', N'عالية جداً', 5, 'very_high', '#DC3545'); -- Red

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM asset_criticality_matrix_factor_ratings WHERE asset_criticality_matrix_id = (SELECT id FROM asset_criticality_matrices WHERE title_en = 'Default Asset Criticality Matrix');
-- +goose StatementEnd
