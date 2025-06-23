-- +goose Up
-- +goose StatementBegin
INSERT INTO process_dependency_sub_types (
    process_dependency_type_id, title_en, title_ar, description_en, description_ar, type, status
) VALUES

-- 🔸 Internal Sub-Types
((SELECT id FROM process_dependency_types WHERE type = 'internal'),
    'Process', 'عملية',
    'A preceding or supporting process within the organization that must be completed or active.',
    'عملية تسبق أو تدعم العملية الحالية داخل المنظمة ويجب إكمالها أو أن تكون نشطة.',
    'process', 'active'),

((SELECT id FROM process_dependency_types WHERE type = 'internal'),
    'Department', 'قسم',
    'A specific internal role, team, or department required to perform, approve, or support the process.',
    'دور أو فريق أو قسم داخلي محدد مطلوب لتنفيذ أو الموافقة أو دعم العملية.',
    'department', 'active'),

((SELECT id FROM process_dependency_types WHERE type = 'internal'),
    'Asset', 'أصل',
    'Software platforms or tools used internally that the process relies on (e.g., ERP, CRM).',
    'منصات أو أدوات برمجية مستخدمة داخليًا تعتمد عليها العملية (مثل ERP أو CRM).',
    'asset', 'active'),

((SELECT id FROM process_dependency_types WHERE type = 'internal'),
    'Other', 'أخرى',
    'Recurring internal events or cycles that trigger or impact the process (e.g., quarterly meetings, monthly reviews).',
    'أحداث أو دورات داخلية متكررة تؤثر أو تحفز العملية (مثل الاجتماعات الفصلية أو المراجعات الشهرية).',
    'other', 'active'),

-- 🔸 External Sub-Types
((SELECT id FROM process_dependency_types WHERE type = 'external'),
    'Vendor / Third Party', 'بائع / طرف ثالث',
    'External service providers or partners whose actions or deliverables affect the process.',
    'مزودو خدمات أو شركاء خارجيون تؤثر أفعالهم أو مخرجاتهم على العملية.',
    'vendor', 'active'),

((SELECT id FROM process_dependency_types WHERE type = 'external'),
    'Regulatory Body', 'هيئة تنظيمية',
    'Government or industry regulators that set conditions, rules, or timelines the process must comply with.',
    'الجهات الحكومية أو التنظيمية التي تضع الشروط أو القواعد أو الجداول الزمنية التي يجب على العملية الالتزام بها.',
    'other', 'active'),

((SELECT id FROM process_dependency_types WHERE type = 'external'),
    'External System / API', 'نظام خارجي / واجهة برمجة التطبيقات',
    'Technology systems or APIs managed outside the organization that are used within the process.',
    'أنظمة تقنية أو واجهات برمجة تطبيقات تُدار خارج المنظمة وتُستخدم داخل العملية.',
    'other', 'active'),

((SELECT id FROM process_dependency_types WHERE type = 'external'),
    'Customer Action', 'إجراء العميل',
    'Activities or approvals required from customers or end users to continue or complete the process.',
    'الأنشطة أو الموافقات المطلوبة من العملاء أو المستخدمين النهائيين لمواصلة أو إكمال العملية.',
    'other', 'active'),

((SELECT id FROM process_dependency_types WHERE type = 'external'),
    'External Event', 'حدث خارجي',
    'External circumstances or scheduled occurrences (e.g., policy changes, market events, environmental triggers) that influence the process.',
    'ظروف خارجية أو أحداث مجدولة (مثل تغييرات السياسة أو الأحداث السوقية أو العوامل البيئية) تؤثر على العملية.',
    'other', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE process_dependency_sub_types;
-- +goose StatementEnd
