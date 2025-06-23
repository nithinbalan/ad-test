-- +goose Up
-- +goose StatementBegin
INSERT INTO meeting_types (title_en, title_ar, status) VALUES
('Compliance Review Meeting', 'اجتماع مراجعة الامتثال', 'active'),
('Risk Review Meeting', 'اجتماع مراجعة المخاطر', 'active'),
('Internal Audit Review Meeting', 'اجتماع مراجعة التدقيق الداخلي', 'active'),
('External Audit Review Meeting', 'اجتماع مراجعة التدقيق الخارجي', 'active'),
('Policy Review Meeting', 'اجتماع مراجعة السياسات', 'active'),
('Procedure Review Meeting', 'اجتماع مراجعة الإجراءات', 'active'),
('Contract Review Meeting', 'اجتماع مراجعة العقود', 'active'),
('SLA Review Meeting', 'اجتماع مراجعة اتفاقية مستوى الخدمة', 'active'),
('Vendor Review Meeting', 'اجتماع مراجعة الموردين', 'active'),
('Performance Review Meeting', 'اجتماع مراجعة الأداء', 'active'),
('Incident Review Meeting', 'اجتماع مراجعة الحوادث', 'active'),
('Strategy Review Meeting', 'اجتماع مراجعة الاستراتيجية', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE meeting_types;
-- +goose StatementEnd
