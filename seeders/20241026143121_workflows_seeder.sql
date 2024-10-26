-- +goose Up
-- Insert seed data for workflows

INSERT INTO workflows (type, status, created_by_id, created_at, module_identifier, title_en, title_ar)
VALUES 
    ('draft-report', 'inactive', 1, GETDATE(), 'audit', 'Audit Draft Report', 'تقرير مسودة التدقيق'),
    ('final-report', 'active', 1, GETDATE(), 'audit', 'Audit Final Report', 'التقرير النهائي للتدقيق');

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE workflows;
-- +goose StatementEnd
