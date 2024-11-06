-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

INSERT INTO asset_criticality_matrices (title_en, title_ar, status)
VALUES ('Default Asset Criticality Matrix', N'مصفوفة أهمية الأصول الافتراضية', 'active');

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM asset_criticality_matrices WHERE title_en = 'Default Asset Criticality Matrix';
-- +goose StatementEnd
