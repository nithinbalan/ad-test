-- +goose Up
-- +goose StatementBegin
INSERT INTO asset_device_categories (title_en, title_ar, type)
VALUES
    ('Servers', N'الخوادم', 'servers'),
    ('Employee devices', N'أجهزة الموظفين', 'employee-devices'),
    ('Databases', N'قواعد البيانات', 'databases'),
    ('Network components', N'مكونات الشبكة', 'network-components'),
    ('Protection systems', N'أنظمة الحماية', 'protection-systems'),
    ('User-management systems', N'أنظمة إدارة المستخدمين', 'user-management-systems'),
    ('Remote-work systems', N'أنظمة العمل عن بُعد', 'remote-work-systems'),
    ('Internet/browsing services', N'خدمات الإنترنت والتصفح', 'internet-browsing-services'),
    ('Printers', N'الطابعات', 'printers'),
    ('Phones', N'الهواتف', 'phones'),
    ('VMs', N'الآلات الافتراضية', 'vms'),
    ('Kubernetes', N'كوبرنيتس', 'kubernetes'),
    ('Routers', N'أجهزة التوجيه', 'routers'),
    ('Switches', N'المبدلات', 'switches');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM asset_device_categories
WHERE type IN (
    'servers',
    'employee-devices',
    'databases',
    'network-components',
    'protection-systems',
    'user-management-systems',
    'remote-work-systems',
    'internet-browsing-services',
    'printers',
    'phones',
    'vms',
    'kubernetes',
    'routers',
    'switches'
);
-- +goose StatementEnd