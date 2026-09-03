-- +goose Up
-- +goose StatementBegin
INSERT INTO planned_advisory_classifications (title_en, title_ar, status) VALUES
    ('General Note', N'ملاحظة عامة', 'active'),
    ('RFI', N'طلب معلومات', 'active'),
    ('Release Information', N'إصدار معلومات', 'active'),
    ('Share Threat Information', N'مشاركة معلومات التهديد', 'active'),
    ('Report Incident', N'الإبلاغ عن حادث', 'active'),
    ('Enquiry', N'استفسار', 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM planned_advisory_classifications
WHERE title_en IN ('General Note', 'RFI', 'Release Information', 'Share Threat Information', 'Report Incident', 'Enquiry');
-- +goose StatementEnd
