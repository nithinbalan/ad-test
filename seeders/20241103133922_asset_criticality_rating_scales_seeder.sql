-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @MatrixId INT;
SELECT @MatrixId = id FROM asset_criticality_matrices WHERE title_en = 'Default Asset Criticality Matrix';

INSERT INTO asset_criticality_rating_scales (asset_criticality_matrix_id, title_en, title_ar, description_en, description_ar, score_from, score_to, type, color_code)
VALUES
    (@MatrixId, 'Low', N'منخفض', 'Assets with low criticality.', N'الأصول ذات الأهمية المنخفضة.', 1.0, 2.0, 'low', '#28A745'), -- Green
    (@MatrixId, 'Medium', N'متوسط', 'Assets with medium criticality.', N'الأصول ذات الأهمية المتوسطة.', 2.1, 3.5, 'medium', '#FFC107'), -- Yellow
    (@MatrixId, 'High', N'عالي', 'Assets with high criticality.', N'الأصول ذات الأهمية العالية.', 3.6, 5.0, 'high', '#DC3545'); -- Red

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM asset_criticality_rating_scales WHERE asset_criticality_matrix_id = (SELECT id FROM asset_criticality_matrices WHERE title_en = 'Default Asset Criticality Matrix');
-- +goose StatementEnd
