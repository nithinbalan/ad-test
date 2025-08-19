-- +goose Up
-- +goose StatementBegin
/* Seed four frameworks (Binary, Point, Percentage, SAMA Maturity) */

DECLARE @typeBinary INT, @typePoint INT, @typePercentage INT, @typeMaturity INT;
SELECT @typeBinary     = id FROM compliance_assessment_framework_types WHERE type = 'binary';
SELECT @typePoint      = id FROM compliance_assessment_framework_types WHERE type = 'point';
SELECT @typePercentage = id FROM compliance_assessment_framework_types WHERE type = 'percentage';
SELECT @typeMaturity   = id FROM compliance_assessment_framework_types WHERE type = 'maturity';

IF @typeBinary IS NULL OR @typePoint IS NULL OR @typePercentage IS NULL OR @typeMaturity IS NULL
    RAISERROR('Framework types missing. Run compliance_assessment_framework_types migration first.', 16, 1);

DECLARE @fwBinary INT, @fwPoint INT, @fwPercentage INT, @fwMaturity INT;

/* Binary */
IF NOT EXISTS (
    SELECT 1 FROM compliance_assessment_frameworks 
    WHERE title_en = 'Binary' AND compliance_assessment_framework_type_id = @typeBinary
)
BEGIN
    INSERT INTO compliance_assessment_frameworks
        (title_en, title_ar, description, compliance_assessment_framework_type_id, status, created_by_id, created_at)
    VALUES
        ('Binary', N'ثنائي', N'Yes/No compliance option with optional NA', @typeBinary, 'active', 1, GETDATE());
    SET @fwBinary = SCOPE_IDENTITY();
END
ELSE
    SELECT @fwBinary = id FROM compliance_assessment_frameworks 
    WHERE title_en = 'Binary' AND compliance_assessment_framework_type_id = @typeBinary;

/* Point Scoring */
IF NOT EXISTS (
    SELECT 1 FROM compliance_assessment_frameworks 
    WHERE title_en = 'Point Scoring' AND compliance_assessment_framework_type_id = @typePoint
)
BEGIN
    INSERT INTO compliance_assessment_frameworks
        (title_en, title_ar, description, compliance_assessment_framework_type_id, status, created_by_id, created_at)
    VALUES
        ('Point Scoring', N'التقييم بالنقاط', 
         N'Scoring by points per requirement (0,1,2) with overall bands', 
         @typePoint, 'active', 1, GETDATE());
    SET @fwPoint = SCOPE_IDENTITY();
END
ELSE
    SELECT @fwPoint = id FROM compliance_assessment_frameworks 
    WHERE title_en = 'Point Scoring' AND compliance_assessment_framework_type_id = @typePoint;

/* Percentage Scoring */
IF NOT EXISTS (
    SELECT 1 FROM compliance_assessment_frameworks 
    WHERE title_en = 'Percentage Scoring' AND compliance_assessment_framework_type_id = @typePercentage
)
BEGIN
    INSERT INTO compliance_assessment_frameworks
        (title_en, title_ar, description, compliance_assessment_framework_type_id, status, created_by_id, created_at)
    VALUES
        ('Percentage Scoring', N'التقييم بالنسبة المئوية', 
         N'Bands: 0–50, 50–80, 80–100 with NA support', 
         @typePercentage, 'active', 1, GETDATE());
    SET @fwPercentage = SCOPE_IDENTITY();
END
ELSE
    SELECT @fwPercentage = id FROM compliance_assessment_frameworks 
    WHERE title_en = 'Percentage Scoring' AND compliance_assessment_framework_type_id = @typePercentage;

-- Replace the old SAMA block with this

IF NOT EXISTS (
    SELECT 1 FROM compliance_assessment_frameworks 
    WHERE title_en = 'Maturity Framework' AND compliance_assessment_framework_type_id = @typeMaturity
)
BEGIN
    INSERT INTO compliance_assessment_frameworks
        (title_en, title_ar, description, compliance_assessment_framework_type_id, status, created_by_id, created_at)
    VALUES
        ('Maturity Framework', N'إطار النضج',
         N'Levels L0–L5 with percentage bands for status mapping',
         @typeMaturity, 'active', 1, GETDATE());
    SET @fwMaturity = SCOPE_IDENTITY();
END
ELSE
    SELECT @fwMaturity = id FROM compliance_assessment_frameworks 
    WHERE title_en = 'Maturity Framework' AND compliance_assessment_framework_type_id = @typeMaturity;

-- (Optional) If you previously seeded "SAMA Maturity Framework", normalize it:
IF @fwMaturity IS NULL
BEGIN
    UPDATE compliance_assessment_frameworks
    SET title_en = 'Maturity Framework', title_ar = N'إطار النضج'
    WHERE title_en = 'SAMA Maturity Framework' AND compliance_assessment_framework_type_id = @typeMaturity;

    SELECT @fwMaturity = id FROM compliance_assessment_frameworks 
    WHERE title_en = 'Maturity Framework' AND compliance_assessment_framework_type_id = @typeMaturity;
END

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE compliance_assessment_frameworks;
-- +goose StatementEnd
