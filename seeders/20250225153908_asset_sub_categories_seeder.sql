-- +goose Up
-- +goose StatementBegin
INSERT INTO asset_sub_categories (asset_category_id, title_en, title_ar, status) VALUES
-- Information Assets
((SELECT id FROM asset_categories WHERE title_en = 'Information Assets'), 'Policies and Procedures', 'السياسات والإجراءات', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Information Assets'), 'Intellectual Property', 'الملكية الفكرية', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Information Assets'), 'Strategic Plans', 'الخطط الاستراتيجية', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Information Assets'), 'Operational Documentation', 'وثائق التشغيل', 'active'),

-- Data Assets
((SELECT id FROM asset_categories WHERE title_en = 'Data Assets'), 'Customer Data', 'بيانات العملاء', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Data Assets'), 'Employee Data', 'بيانات الموظفين', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Data Assets'), 'Financial Records', 'السجلات المالية', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Data Assets'), 'Audit Logs', 'سجلات التدقيق', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Data Assets'), 'Research Data', 'بيانات البحث', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Data Assets'), 'Sensitive/Confidential Data', 'البيانات الحساسة/السرية', 'active'),

-- Business Applications
((SELECT id FROM asset_categories WHERE title_en = 'Business Applications'), 'Enterprise Resource Planning (ERP) Systems', 'أنظمة تخطيط موارد المؤسسات', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Business Applications'), 'Customer Relationship Management (CRM) Software', 'برامج إدارة علاقات العملاء', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Business Applications'), 'Business Intelligence (BI) Tools', 'أدوات ذكاء الأعمال', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Business Applications'), 'Office Productivity Applications', 'تطبيقات الإنتاجية المكتبية', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Business Applications'), 'Custom Applications', 'التطبيقات المخصصة', 'active'),

-- Hardware Assets
((SELECT id FROM asset_categories WHERE title_en = 'Hardware Assets'), 'Servers (Physical and Virtual)', 'الخوادم (المادية والافتراضية)', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Hardware Assets'), 'Workstations (Desktops and Laptops)', 'محطات العمل (أجهزة الكمبيوتر المكتبية والمحمولة)', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Hardware Assets'), 'Mobile Devices (Tablets and Smartphones)', 'الأجهزة المحمولة (الأجهزة اللوحية والهواتف الذكية)', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Hardware Assets'), 'Networking Devices', 'أجهزة الشبكة', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Hardware Assets'), 'IoT Devices', 'أجهزة إنترنت الأشياء', 'active'),

-- Software Assets
((SELECT id FROM asset_categories WHERE title_en = 'Software Assets'), 'Operating Systems', 'أنظمة التشغيل', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Software Assets'), 'Licensed Software', 'البرامج المرخصة', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Software Assets'), 'Development Tools', 'أدوات التطوير', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Software Assets'), 'Security Tools', 'أدوات الأمان', 'active'),

-- Cloud Assets
((SELECT id FROM asset_categories WHERE title_en = 'Cloud Assets'), 'Cloud Infrastructure', 'البنية التحتية السحابية', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Cloud Assets'), 'SaaS Applications', 'تطبيقات SaaS', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Cloud Assets'), 'Storage', 'التخزين', 'active'),

-- Network Assets
((SELECT id FROM asset_categories WHERE title_en = 'Network Assets'), 'Network Segments', 'أجزاء الشبكة', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Network Assets'), 'Virtual Private Networks (VPN)', 'الشبكات الخاصة الافتراضية', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Network Assets'), 'Internet Gateways', 'بوابات الإنترنت', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Network Assets'), 'Intrusion Detection Systems (IDS)', 'أنظمة كشف التسلل', 'active'),

-- People Assets
((SELECT id FROM asset_categories WHERE title_en = 'People Assets'), 'Employees', 'الموظفون', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'People Assets'), 'Contractors', 'المقاولون', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'People Assets'), 'Third-party Vendors', 'البائعون الخارجيون', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'People Assets'), 'Administrators (Privileged Users)', 'المسؤولون (المستخدمون المتميزون)', 'active'),

-- Physical Assets
((SELECT id FROM asset_categories WHERE title_en = 'Physical Assets'), 'Office Buildings', 'المباني المكتبية', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Physical Assets'), 'Data Centers', 'مراكز البيانات', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Physical Assets'), 'Backup Storage Facilities', 'منشآت تخزين النسخ الاحتياطية', 'active'),

-- Communication Assets
((SELECT id FROM asset_categories WHERE title_en = 'Communication Assets'), 'Email Systems', 'أنظمة البريد الإلكتروني', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Communication Assets'), 'Collaboration Tools', 'أدوات التعاون', 'active'),
((SELECT id FROM asset_categories WHERE title_en = 'Communication Assets'), 'Telecommunication Systems', 'أنظمة الاتصالات', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM asset_sub_categories WHERE asset_category_id IN (SELECT id FROM asset_categories WHERE title_en IN (
'Information Assets', 'Data Assets', 'Business Applications', 
'Hardware Assets', 'Software Assets', 'Cloud Assets', 
'Network Assets', 'People Assets', 'Physical Assets', 
'Communication Assets'));
-- +goose StatementEnd
