-- +goose Up
-- +goose StatementBegin

DECLARE @fwBinary INT, @fwPoint INT, @fwPercentage INT, @fwMaturity INT;
SELECT @fwBinary     = id FROM compliance_assessment_frameworks WHERE title_en = 'Binary';
SELECT @fwPoint      = id FROM compliance_assessment_frameworks WHERE title_en = 'Point Scoring';
SELECT @fwPercentage = id FROM compliance_assessment_frameworks WHERE title_en = 'Percentage Scoring';
SELECT @fwMaturity   = id FROM compliance_assessment_frameworks WHERE title_en = 'Maturity Framework';

IF @fwBinary IS NULL OR @fwPoint IS NULL OR @fwPercentage IS NULL OR @fwMaturity IS NULL
    RAISERROR('Frameworks missing. Run frameworks seeder first.', 16, 1);

DECLARE @stCompliant INT, @stPartial INT, @stNotImpl INT, @stNA INT;
SELECT @stCompliant = id FROM compliance_statuses WHERE type = 'compliant';
SELECT @stPartial   = id FROM compliance_statuses WHERE type = 'partially_compliant';
SELECT @stNotImpl   = id FROM compliance_statuses WHERE type = 'not_implemented';
SELECT @stNA        = id FROM compliance_statuses WHERE type = 'not_applicable';

IF @stCompliant IS NULL OR @stPartial IS NULL OR @stNotImpl IS NULL OR @stNA IS NULL
    RAISERROR('Compliance statuses missing.', 16, 1);

/* ===== Percentage (three ranges; value is NULL, min/max REQUIRED) ===== */
IF NOT EXISTS (SELECT 1 FROM compliance_assessment_framework_values WHERE compliance_assessment_framework_id = @fwPercentage)
BEGIN
    INSERT INTO compliance_assessment_framework_values
        (compliance_assessment_framework_id, compliance_status_id, value, title_en, title_ar, description,
         score_min, score_max, is_applicable, created_by_id, created_at)
    VALUES
        (@fwPercentage, @stNotImpl,  NULL, 'Not Implemented',        N'غير مُنفّذ',            N'0–50%',     0,   50, 1, 1, GETDATE()),
        (@fwPercentage, @stPartial,  NULL, 'Partially Implemented',  N'تم التنفيذ جزئياً',     N'50–80%',   50,   80, 1, 1, GETDATE()),
        (@fwPercentage, @stCompliant,NULL, 'Compliant',              N'متوافق',                 N'80–100%',  80,  100, 1, 1, GETDATE()),
         (@fwPercentage, @stNA, NULL, 'Not Applicable', N'غير قابل للتطبيق', N'NA', NULL, NULL, 0, 1, GETDATE());
END

/* ===== Point (0/1/2) each WITH min/max (no NA, no extra bands) ===== */
IF NOT EXISTS (SELECT 1 FROM compliance_assessment_framework_values WHERE compliance_assessment_framework_id = @fwPoint)
BEGIN
    INSERT INTO compliance_assessment_framework_values
        (compliance_assessment_framework_id, compliance_status_id, value, title_en, title_ar, description,
         score_min, score_max, is_applicable, created_by_id, created_at)
    VALUES
        (@fwPoint, @stNotImpl,  0, '0 Point', N'0 نقطة',  N'Not implemented',        0,   50, 1, 1, GETDATE()),
        (@fwPoint, @stPartial,  1, '1 Point', N'1 نقطة',  N'Partially implemented',  50,  80, 1, 1, GETDATE()),
        (@fwPoint, @stCompliant,2, '2 Points',N'2 نقطتان',N'Fully implemented',      80, 100, 1, 1, GETDATE()),
        (@fwPoint, @stNA, -1, 'Not Applicable', N'غير قابل للتطبيق', N'NA', NULL, NULL, 0, 1, GETDATE());
END

/* ===== Binary (Yes/No) WITH min/max (no NA, no extra bands) ===== */
IF NOT EXISTS (SELECT 1 FROM compliance_assessment_framework_values WHERE compliance_assessment_framework_id = @fwBinary)
BEGIN
    INSERT INTO compliance_assessment_framework_values
        (compliance_assessment_framework_id, compliance_status_id, value, title_en, title_ar, description,
         score_min, score_max, is_applicable, created_by_id, created_at)
    VALUES
        (@fwBinary, @stNotImpl,   0, 'No',  N'لا',  N'Not implemented',  0,   50, 1, 1, GETDATE()),
        (@fwBinary, @stCompliant, 1, 'Yes', N'نعم', N'Implemented',      50, 100, 1, 1, GETDATE()),
        (@fwBinary, @stNA, 0, 'Not Applicable', N'غير قابل للتطبيق', N'NA', NULL, NULL, 0, 1, GETDATE());
END

/* ===== Maturity (L0–L5) each WITH min/max (no NA) ===== */
IF NOT EXISTS (SELECT 1 FROM compliance_assessment_framework_values WHERE compliance_assessment_framework_id = @fwMaturity)
BEGIN
    INSERT INTO compliance_assessment_framework_values
        (compliance_assessment_framework_id, compliance_status_id, value, title_en, title_ar, description,
         score_min, score_max, is_applicable, created_by_id, created_at)
    VALUES
        (@fwMaturity, @stNotImpl, 0, 'Level 0 — Non-existent',       N'المستوى 0 — غير موجود',           N'No documentation or awareness',                  0,   20, 1, 1, GETDATE()),
        (@fwMaturity, @stNotImpl, 1, 'Level 1 — Ad-hoc',             N'المستوى 1 — عشوائي',               N'Informal/occasional; inconsistent practice',    20,  40, 1, 1, GETDATE()),
        (@fwMaturity, @stPartial, 2, 'Level 2 — Repeatable',         N'المستوى 2 — متكرر',                N'Repeatable but informal',                        40,  60, 1, 1, GETDATE()),
        (@fwMaturity, @stPartial, 3, 'Level 3 — Structured',         N'المستوى 3 — منظم ورَسمي',           N'Defined and implemented formally',               60,  80, 1, 1, GETDATE()),
        (@fwMaturity, @stCompliant,4,'Level 4 — Managed & Measured', N'المستوى 4 — مُدار ومقاس',           N'Measured and performance managed',               80,  95, 1, 1, GETDATE()),
        (@fwMaturity, @stCompliant,5,'Level 5 — Adaptive',           N'المستوى 5 — تكيفي',                 N'Continuous improvement',                         95, 100, 1, 1, GETDATE()),
        (@fwMaturity, @stNA, -1, 'Not Applicable', N'غير قابل للتطبيق', N'NA', NULL, NULL, 0, 1, GETDATE());
END
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE compliance_assessment_framework_values;
-- +goose StatementEnd
