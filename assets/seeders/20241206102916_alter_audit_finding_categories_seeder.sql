-- +goose Up
-- +goose StatementBegin
UPDATE audit_finding_categories
SET title_en = CASE 
                  WHEN title_en = 'Major Non-Compliance' THEN 'Major Non-Conformity (NC)'
                  WHEN title_en = 'Minor Non-Compliance' THEN 'Minor Non-Conformity (NC)'
                  WHEN title_en = 'Opportunity for Improvement' THEN 'Opportunity for Improvement (OFI)'
                  WHEN title_en = 'Observation' THEN 'Observation'
               END,
    title_ar = CASE
                  WHEN title_ar = 'عدم الامتثال الرئيسي' THEN 'عدم المطابقة الرئيسية'
                  WHEN title_ar = 'عدم الامتثال البسيط' THEN 'عدم المطابقة البسيطة'
                  WHEN title_ar = 'فرصة للتحسين' THEN 'فرصة للتحسين'
                  WHEN title_ar = 'ملاحظة' THEN 'ملاحظة'
               END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE audit_finding_categories
SET title_en = CASE 
                  WHEN title_en = 'Major Non-Conformity (NC)' THEN 'Major Non-Compliance'
                  WHEN title_en = 'Minor Non-Conformity (NC)' THEN 'Minor Non-Compliance'
                  WHEN title_en = 'Opportunity for Improvement (OFI)' THEN 'Opportunity for Improvement'
                  WHEN title_en = 'Observation' THEN 'Observation'
               END,
    title_ar = CASE
                  WHEN title_ar = 'عدم المطابقة الرئيسية' THEN 'عدم الامتثال الرئيسي'
                  WHEN title_ar = 'عدم المطابقة البسيطة' THEN 'عدم الامتثال البسيط'
                  WHEN title_ar = 'فرصة للتحسين' THEN 'فرصة للتحسين'
                  WHEN title_ar = 'ملاحظة' THEN 'ملاحظة'
               END;
-- +goose StatementEnd
