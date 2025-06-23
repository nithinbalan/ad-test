-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @MatrixId INT;
SELECT @MatrixId = id FROM asset_criticality_matrices WHERE title_en = 'Default Asset Criticality Matrix';

INSERT INTO asset_criticality_matrix_factors (asset_criticality_matrix_id, title_en, title_ar, description_en, description_ar, weightage)
VALUES
    (@MatrixId, 'Confidentiality', N'سرية', 'Importance of keeping asset data secure.', N'أهمية الحفاظ على أمان بيانات الأصل.', 0.3),
    (@MatrixId, 'Integrity', N'تكامل', 'Ensuring the accuracy and reliability of data.', N'ضمان دقة وموثوقية البيانات.', 0.2),
    (@MatrixId, 'Availability', N'التوافر', 'Requirement for the asset to be operational and accessible.', N'متطلبات أن يكون الأصل عمليًا ومتوافرًا.', 0.25),
    (@MatrixId, 'Business Impact', N'تأثير العمل', 'Impact on the organization if the asset is compromised or unavailable.', N'التأثير على المؤسسة في حالة تعرض الأصل للخطر أو عدم توفره.', 0.15),
    (@MatrixId, 'Dependency', N'اعتماد', 'Number and importance of other systems that depend on this asset.', N'عدد وأهمية الأنظمة الأخرى التي تعتمد على هذا الأصل.', 0.1);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM asset_criticality_matrix_factors WHERE asset_criticality_matrix_id = (SELECT id FROM asset_criticality_matrices WHERE title_en = 'Default Asset Criticality Matrix');
-- +goose StatementEnd
