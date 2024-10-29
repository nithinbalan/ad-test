-- +goose Up
-- Insert seed data for workflows

INSERT INTO workflows (type, status, created_by_id, created_at, module_identifier, title_en, title_ar)
VALUES 
    ('final-report', 'active', 1, GETDATE(), 'audit', 'Audit Final Report', 'التقرير النهائي للتدقيق');

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE workflows;
-- +goose StatementEnd
