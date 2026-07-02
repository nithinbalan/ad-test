-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @MatrixId INT;
SELECT @MatrixId = id FROM asset_criticality_matrices WHERE title_en = 'Default Asset Criticality Matrix' AND status = 'active';

INSERT INTO asset_criticality_matrix_factor_ratings (asset_criticality_matrix_id, title_en, title_ar, score, type, color_code)
VALUES
    (@MatrixId, 'N/A', N'غير منطبق', 0, 'na', '#6C757D'), 
    (@MatrixId, 'Low', N'منخفض', 1, 'low', '#28A745'), 
    (@MatrixId, 'Medium', N'متوسط', 2, 'medium', '#FFC107'), 
    (@MatrixId, 'High', N'عالي', 3, 'high', '#FFA500'), 
    (@MatrixId, 'Very High', N'عالية جداً', 4, 'very_high', '#DC3545'); 

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @MatrixId INT;
SELECT @MatrixId = id FROM asset_criticality_matrices WHERE title_en = 'Default Asset Criticality Matrix' AND status = 'active';

DELETE FROM asset_criticality_matrix_factor_ratings WHERE asset_criticality_matrix_id = @MatrixId;

COMMIT;
-- +goose StatementEnd