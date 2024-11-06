-- +goose Up
-- +goose StatementBegin
INSERT INTO integrations (integration_category_id,title_en, title_ar, description_en, description_ar,logo_url,api_base_url,is_enable, created_at, created_by_id) VALUES
(1,'LDAP', 'لداب', 'LDAP is a protocol for managing directory services over a network.', 'LDAP هو بروتوكول لإدارة خدمات الدليل عبر الشبكة.','LDAP',NULL, 0,GETDATE(),1),
(1,'Azure AD', 'أزور م', 'Active Directory not built for remote work-try JumpCloud free.', 'Active Directory لم يتم تصميمه للعمل عن بعد - جرب JumpCloud مجانًا.','Azure AD',NULL, 0,GETDATE(),1),
(1,'Saudi National ID', 'الهوية الوطنية السعودية', 'The Saudi National ID Card is proof of identity and citizenship for Saudi citizens.', 'بطاقة الهوية الوطنية السعودية هي إثبات الهوية والجنسية للمواطنين السعوديين.','Saudi National ID',NULL, 0,GETDATE(),1),
(2,'SharePoint', 'شيربوينت', 'Share and Manage content for teamwork, easy access, and collaboration.', 'مشاركة وإدارة المحتوى للعمل الجماعي وسهولة الوصول والتعاون.','SharePoint',NULL, 0,GETDATE(),1),
(2,'OneDrive', 'ون درايف', 'Save your files and photos to OneDrive and access them from any device, anywhere.', 'احفظ ملفاتك وصورك في OneDrive وقم بالوصول إليها من أي جهاز وفي أي مكان.','OneDrive',NULL, 0,GETDATE(),1),
(2,'Google Drive', 'جوجل درايف', 'Google Drive lets you securely store, organize, and collaborate.', 'يتيح لك Google Drive التخزين والتنظيم والتعاون بشكل آمن.','Google Drive',NULL, 0,GETDATE(),1),
(3,'Tenable', 'يمكن الدفاع عنه', 'Tenable helps you find, prioritize, and fix cyber risks with security tools.', 'يساعدك Tenable على العثور على المخاطر السيبرانية وتحديد أولوياتها وإصلاحها باستخدام أدوات الأمان.','Tenable',NULL, 0,GETDATE(),1),
(3,'Qualys', 'كواليس', 'Streamline Risk Management Across Your Enterprise.', 'تبسيط إدارة المخاطر عبر مؤسستك.','Qualys',NULL, 0,GETDATE(),1),
(4,'ManageEngine', 'ManageEngine', 'Management with UEM: automate tasks, manage devices and enhance security.', 'الإدارة باستخدام UEM: أتمتة المهام وإدارة الأجهزة وتعزيز الأمان.','ManageEngine',NULL, 0,GETDATE(),1),
(4,'IBM', 'آي بي إم', 'IBM Business Process Manager unifies human and integration capabilities.', 'يقوم IBM Business Process Manager بتوحيد القدرات البشرية والتكاملية.','IBM',NULL, 0,GETDATE(),1),
(4,'Oracle', 'أوراكل', 'Oracle provides a fully integrated stack of cloud applications and platform services.', 'توفر Oracle مجموعة متكاملة تمامًا من التطبيقات السحابية وخدمات النظام الأساسي.','Oracle',NULL, 0,GETDATE(),1),
(5,'Odoo', 'أودو', 'Odoo offers an integrated suite of business applications to manage operations.', 'يقدم Odoo مجموعة متكاملة من تطبيقات الأعمال لإدارة العمليات.','Odoo',NULL, 0,GETDATE(),1),
(6,'Jira', 'جيرا', 'Affordable management solutions for all. Explore "per user" pricing.', 'حلول إدارية ميسورة التكلفة للجميع. استكشف التسعير لكل مستخدم.','Jira',NULL, 0,GETDATE(),1),
(7,'ChatGPT', 'ChatGPT', 'ChatGPT helps you get answers, find inspiration and be more productive.', 'يساعدك ChatGPT في الحصول على الإجابات والعثور على الإلهام وزيادة إنتاجيتك.','ChatGPT','https://ai-assistant-dev-api.jethur.com/rewrite', 0,GETDATE(),1),
(7,'Gemini', 'تَوأَم', 'Get help with writing, planning, learn, and more from Google AI.', 'احصل على مساعدة في الكتابة والتخطيط والتعلم والمزيد من Google AI.','Gemini',NULL, 0,GETDATE(),1),
(7,'Notion', 'فكرة', 'An all-in-one process kit to help you design logos with more clarity and confidence.', 'مجموعة أدوات معالجة شاملة لمساعدتك على تصميم الشعارات بمزيد من الوضوح والثقة.','Notion',NULL, 0,GETDATE(),1);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE integrations;
-- +goose StatementEnd
