-- +goose Up
-- +goose StatementBegin
INSERT INTO process_dependency_sub_types (
    process_dependency_type_id, title_en, title_ar, description_en, description_ar, status
) VALUES
-- Internal Dependencies
((SELECT id FROM process_dependency_types WHERE type = 'internal'),
    'Technology', 'تكنولوجيا', 'ERP System', 'نظام تخطيط موارد المؤسسة', 'active'),
((SELECT id FROM process_dependency_types WHERE type = 'internal'),
    'People', 'أشخاص', 'Warehouse Staff', 'موظفو المستودع', 'active'),
((SELECT id FROM process_dependency_types WHERE type = 'internal'),
    'Operations', 'العمليات', 'Production Machinery', 'آلات الإنتاج', 'active'),
((SELECT id FROM process_dependency_types WHERE type = 'internal'),
    'Infrastructure', 'البنية التحتية', 'Office Network', 'شبكة المكتب', 'active'),

-- External Dependencies
((SELECT id FROM process_dependency_types WHERE type = 'external'),
    'Vendor', 'بائع', 'Logistics Partner', 'شريك لوجستي', 'active'),
((SELECT id FROM process_dependency_types WHERE type = 'external'),
    'Service Provider', 'مزود خدمة', 'Payment Gateway', 'بوابة الدفع', 'active'),
((SELECT id FROM process_dependency_types WHERE type = 'external'),
    'Supplier', 'المورد', 'Raw Material Supplier', 'مورد المواد الخام', 'active'),
((SELECT id FROM process_dependency_types WHERE type = 'external'),
    'Consultant', 'استشاري', 'IT Security Consultant', 'استشاري أمن تكنولوجيا المعلومات', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE process_dependency_sub_types;
-- +goose StatementEnd
