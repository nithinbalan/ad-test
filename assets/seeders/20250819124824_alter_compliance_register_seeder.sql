-- +goose Up
-- +goose StatementBegin
DECLARE @fwMaturity INT;

SELECT @fwMaturity = f.id
FROM compliance_assessment_frameworks f
JOIN compliance_assessment_framework_types t
  ON f.compliance_assessment_framework_type_id = t.id
WHERE t.type = 'maturity';

IF @fwMaturity IS NULL
    RAISERROR('Maturity Framework not found. Run frameworks seeder first.', 16, 1);

-- Update compliance registers that don’t have a framework assigned
UPDATE compliance_register
SET compliance_assessment_framework_id = @fwMaturity
WHERE compliance_assessment_framework_id IS NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Rollback: reset compliance_register that were set to Maturity back to NULL
DECLARE @fwMaturity INT;

SELECT @fwMaturity = f.id
FROM compliance_assessment_frameworks f
JOIN compliance_assessment_framework_types t
  ON f.compliance_assessment_framework_type_id = t.id
WHERE t.type = 'maturity';

UPDATE compliance_register
SET compliance_assessment_framework_id = NULL
WHERE compliance_assessment_framework_id = @fwMaturity;
-- +goose StatementEnd
