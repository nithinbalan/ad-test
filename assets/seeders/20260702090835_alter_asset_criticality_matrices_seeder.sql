-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

UPDATE asset_criticality_matrices
SET
    title_en = 'Default Asset Criticality Matrix - Inactive',
    title_ar = N'مصفوفة أهمية الأصول الافتراضية - غير نشط',
    status = 'inactive'
WHERE title_en = 'Default Asset Criticality Matrix';

INSERT INTO asset_criticality_matrices (
    title_en,
    title_ar,
    status
)
VALUES (
    'Default Asset Criticality Matrix',
    N'مصفوفة أهمية الأصول الافتراضية',
    'active'
);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN TRANSACTION;

DELETE FROM asset_criticality_matrices
WHERE title_en = 'Default Asset Criticality Matrix'
  AND status = 'active';

UPDATE asset_criticality_matrices
SET
    title_en = 'Default Asset Criticality Matrix',
    title_ar = N'مصفوفة أهمية الأصول الافتراضية',
    status = 'active'
WHERE title_en = 'Default Asset Criticality Matrix - Inactive';

COMMIT;
-- +goose StatementEnd