-- +goose Up
-- +goose StatementBegin
INSERT INTO asset_device_categories (title_en, title_ar, is_system)
VALUES
    ('Servers', N'الخوادم', 1),
    ('Employee devices', N'أجهزة الموظفين', 1),
    ('Databases', N'قواعد البيانات', 1),
    ('Network components', N'مكونات الشبكة', 1),
    ('Protection systems', N'أنظمة الحماية', 1),
    ('User-management systems', N'أنظمة إدارة المستخدمين', 1),
    ('Remote-work systems', N'أنظمة العمل عن بُعد', 1),
    ('Internet/browsing services', N'خدمات الإنترنت والتصفح', 1),
    ('Printers', N'الطابعات', 1),
    ('Phones', N'الهواتف', 1),
    ('VMs', N'الآلات الافتراضية', 1),
    ('Kubernetes', N'كوبرنيتس', 1),
    ('Routers', N'أجهزة التوجيه', 1),
    ('Switches', N'المبدلات', 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM asset_device_categories
WHERE is_system = 1;
-- +goose StatementEnd