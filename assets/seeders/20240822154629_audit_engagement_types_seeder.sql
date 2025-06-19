-- +goose Up
INSERT INTO audit_engagement_types (title_en, title_ar, description_en, description_ar)
VALUES 
    ('Compliance Verification Meeting', 'اجتماع التحقق من الامتثال', 'A meeting to ensure that all regulatory requirements are being met during the audit process.', 'اجتماع لضمان تلبية جميع المتطلبات التنظيمية خلال عملية التدقيق.'),
    ('Audit Exit Meeting', 'اجتماع اختتام التدقيق', 'A final meeting to review audit findings and recommendations with senior management and key stakeholders.', 'اجتماع نهائي لمراجعة نتائج التدقيق والتوصيات مع الإدارة العليا وأصحاب المصلحة الرئيسيين.'),
    ('Risk Discussion Session', 'جلسة مناقشة المخاطر', 'A session to discuss identified risks and determine appropriate mitigation strategies.', 'جلسة لمناقشة المخاطر التي تم تحديدها وتحديد استراتيجيات التخفيف المناسبة.'),
    ('Process Walkthrough', 'مراجعة العملية', 'A detailed walkthrough of specific processes with process owners to identify potential weaknesses.', 'مراجعة تفصيلية لعمليات محددة مع أصحاب العمليات لتحديد نقاط الضعف المحتملة.'),
    ('Audit Report Presentation', 'عرض تقرير التدقيق', 'A presentation of the final audit report, explaining findings, risks, and suggested improvements.', 'عرض التقرير النهائي للتدقيق، يتضمن شرح النتائج والمخاطر والتحسينات المقترحة.');

-- +goose Down
TRUNCATE TABLE audit_engagement_types;
