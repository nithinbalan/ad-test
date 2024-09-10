-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_activity_types
    (title_en, title_ar, description_en, description_ar, status)
VALUES
    ('Pre-Audit Meeting', 'اجتماع ما قبل التدقيق', 'Meeting before the audit begins', 'اجتماع قبل بدء التدقيق', 'active'),
    ('Audit Execution', 'تنفيذ التدقيق', 'Execution of the audit process', 'تنفيذ عملية التدقيق', 'active'),
    ('Follow-Up Actions', 'إجراءات المتابعة', 'Actions taken after the audit', 'الإجراءات المتخذة بعد التدقيق', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_activity_types;
-- +goose StatementEnd