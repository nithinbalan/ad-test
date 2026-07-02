-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @MatrixId INT;

SELECT
    @MatrixId = id
FROM
    asset_criticality_matrices
WHERE
    title_en = 'Default Asset Criticality Matrix'
    AND status = 'active';

DELETE FROM asset_criticality_matrix_factors
WHERE asset_criticality_matrix_id = @MatrixId;

INSERT INTO asset_criticality_matrix_factors (
    asset_criticality_matrix_id,
    title_en,
    title_ar,
    description_en,
    description_ar,
    weightage
)
VALUES
(
    @MatrixId,
    'Confidentiality',
    N'السرية',
    'Importance of keeping asset data secure.',
    'أهمية الحفاظ على أمان بيانات الأصل.',
    0.33
),
(
    @MatrixId,
    'Integrity',
    N'التكامل',
    'Ensuring the accuracy and reliability of data.',
    N'ضمان دقة وموثوقية البيانات.',
    0.33
),
(
    @MatrixId,
    'Availability',
    N'التوافر',
    'Requirement for the asset to be operational and accessible.',
    N'متطلبات أن يكون الأصل عمليًا ومتاحًا.',
    0.33
);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @MatrixId INT;

SELECT
    @MatrixId = id
FROM
    asset_criticality_matrices
WHERE
    title_en = 'Default Asset Criticality Matrix'
    AND status = 'active';

DELETE FROM asset_criticality_matrix_factors
WHERE asset_criticality_matrix_id = @MatrixId;

COMMIT;
-- +goose StatementEnd