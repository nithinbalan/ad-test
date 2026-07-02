-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @MatrixId INT;
SELECT @MatrixId = id FROM asset_criticality_matrices WHERE title_en = 'Default Asset Criticality Matrix' AND status = 'active';

INSERT INTO asset_criticality_matrix_rating_scales (asset_criticality_matrix_id, title_en, title_ar, description_en, description_ar, score_from, score_to, type, color_code)
VALUES
    (@MatrixId, 'Low', N'منخفض', 'Assets with low criticality.', N'الأصول ذات الأهمية المنخفضة.', 1.0, 2.0, 'low', '#28A745'), 
    (@MatrixId, 'Medium', N'متوسط', 'Assets with medium criticality.', N'الأصول ذات الأهمية المتوسطة.', 2.0, 3.0, 'medium', '#FFC107'), 
    (@MatrixId, 'High', N'عالي', 'Assets with high criticality.', N'الأصول ذات الأهمية العالية.', 3.0, 4.0, 'high', '#FFA500'), 
    (@MatrixId, 'Critical', N'حرج', 'Assets with critical criticality.', N'الأصول ذات الأهمية الحرجة.', 4.0, 5.0, 'critical', '#DC3545'); 

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @MatrixId INT;
SELECT @MatrixId = id FROM asset_criticality_matrices WHERE title_en = 'Default Asset Criticality Matrix' AND status = 'active';

DELETE FROM asset_criticality_matrix_rating_scales WHERE asset_criticality_matrix_id = @MatrixId;

COMMIT;
-- +goose StatementEnd