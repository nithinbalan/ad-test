-- +goose Up
-- +goose StatementBegin
INSERT INTO integration_categories (title_en, title_ar, description_ar, description_en, status) VALUES
('Single sign-on', 'تسجيل الدخول الموحد', 'نظام لتسجيل الدخول بكلمة مرور واحدة للوصول إلى جميع الأنظمة', 'A system for logging in with a single password to access all systems', 'active'),
('Document Management', 'إدارة الوثائق', 'نظام لتنظيم وتخزين وإدارة الوثائق الإلكترونية والملفات', 'A system to organize, store, and manage electronic documents and files', 'active'),
('Vulnerability Management', 'إدارة الثغرات الأمنية', 'نظام لتحديد وإدارة وإصلاح الثغرات الأمنية في الأنظمة', 'A system to identify, manage, and remediate security vulnerabilities in systems', 'active'),
('Assets', 'أصول', 'نظام لتتبع وإدارة الأصول والممتلكات التابعة للشركة', 'A system for tracking and managing company assets and properties', 'active'),
('Processes', 'العمليات', 'إدارة العمليات والإجراءات لتحسين الكفاءة والإنتاجية', 'Managing processes and workflows to enhance efficiency and productivity', 'active'),
('Project Management', 'إدارة المشاريع', 'نظام لتخطيط وتنفيذ المشاريع وتتبع التقدم والأداء', 'A system for planning and executing projects while tracking progress and performance', 'active'),
('Artificial Intelligence', 'الذكاء الاصطناعي', 'تطبيق تقنيات الذكاء الاصطناعي لتحسين سير العمل والقرارات', 'Applying AI technologies to improve workflows and decision-making', 'active'),
('Control', 'يتحكم', 'نظام للرقابة والتحكم في العمليات والأنظمة', 'A system for monitoring and controlling processes and systems', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE integration_categories;
-- +goose StatementEnd
