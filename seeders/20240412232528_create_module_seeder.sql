-- +goose Up
-- +goose StatementBegin
DECLARE @orgModuleID INT;

INSERT INTO modules (identifier, title_en, title_ar, [order])
VALUES ('organization', 'Organization', 'المنظمة', 1);

SELECT @orgModuleID = SCOPE_IDENTITY();

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order]) VALUES
('strategic_profile', 'Strategic Profile', 'الملف الاستراتيجي', @orgModuleID, 2),
('user_management', 'User Management', 'إدارة المستخدمين', @orgModuleID, 3),
('process', 'Process', 'العمليات', @orgModuleID, 4),
('controls', 'Controls', 'ضوابط', @orgModuleID, 5);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE modules;
-- +goose StatementEnd
