-- +goose Up
-- +goose StatementBegin

-- Update 'Active Directory' integration
UPDATE integrations
SET 
    description_en = 'Provides secure single sign-on (SSO) and automatic user synchronization using your on-premises Active Directory. Manage authentication and user accounts from a centralized directory.',
    description_ar = 'يوفر تسجيل دخول موحد آمن (SSO) ومزامنة تلقائية للمستخدم باستخدام Active Directory المحلي الخاص بك. قم بإدارة المصادقة وحسابات المستخدم من دليل مركزي.'
WHERE type = 'active-directory';

-- Update 'SIEM' integration
UPDATE integrations
SET 
    title_en = 'SIEM Log Forwarding',
    title_ar = 'إعادة توجيه سجلات SIEM',
    description_en = 'Automatically forward security events and activity logs to your SIEM platform for centralized analysis, correlation, and threat detection.',
    description_ar = 'قم بإعادة توجيه الأحداث الأمنية وسجلات النشاط تلقائيًا إلى منصة SIEM الخاصة بك للتحليل المركزي والربط واكتشاف التهديدات.'
WHERE type = 'siem';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert 'Active Directory' integration description
UPDATE integrations
SET 
    description_en = 'Active Directory not built for remote work-try JumpCloud free.',
    description_ar = 'Active Directory لم يتم تصميمه للعمل عن بعد - جرب JumpCloud مجانًا.'
WHERE type = 'active-directory';

-- Revert 'SIEM Log Forwarding' back to 'SIEM'
UPDATE integrations
SET 
    title_en = 'SIEM',
    title_ar = 'سييم',
    description_en = 'Security Information and Event Management integration for forwarding logs to external systems.',
    description_ar = 'تكامل إدارة معلومات وأحداث الأمان لإرسال السجلات إلى الأنظمة الخارجية.'
WHERE type = 'siem';

-- +goose StatementEnd
