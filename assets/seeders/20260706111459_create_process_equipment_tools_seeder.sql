-- +goose Up
-- +goose StatementBegin
INSERT INTO process_equipment_tools (title_en, title_ar, created_by_id, created_at, status) VALUES
    ('Desktop Computer', N'كمبيوتر مكتبي', 1, GETDATE(), 'active'),
    ('Desk Phone', N'هاتف مكتبي', 1, GETDATE(), 'active'),
    ('Printer', N'طابعة', 1, GETDATE(), 'active'),
    ('Scanner', N'ماسح ضوئي', 1, GETDATE(), 'active'),
    ('Desk / Workspace', N'مكتب / مساحة عمل', 1, GETDATE(), 'active'),
    ('Portable Laptop', N'حاسوب محمول', 1, GETDATE(), 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM process_equipment_tools;
-- +goose StatementEnd
