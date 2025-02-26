-- +goose Up
-- +goose StatementBegin

INSERT INTO asset_categories (title_en, title_ar, status)
VALUES
    ('Hardware', 'الأجهزة', 'active'),
    ('Software', 'البرمجيات', 'active'),
    ('Networking Equipment', 'معدات الشبكات', 'active'),
    ('Cloud Resources', 'الموارد السحابية', 'active'),
    ('Buildings', 'المباني', 'active'),
    ('Vehicles', 'المركبات', 'active'),
    ('Machinery', 'الآلات', 'active'),
    ('Furniture', 'الأثاث', 'active'),
    ('Employees', 'الموظفين', 'active'),
    ('Contractors', 'المقاولون', 'active'),
    ('Training Programs', 'برامج التدريب', 'active'),
    ('Recruitment Resources', 'موارد التوظيف', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM asset_categories WHERE title_en IN (
    'Hardware', 'Software', 'Networking Equipment', 'Cloud Resources',
    'Buildings', 'Vehicles', 'Machinery', 'Furniture',
    'Employees', 'Contractors', 'Training Programs', 'Recruitment Resources'
);

-- +goose StatementEnd
