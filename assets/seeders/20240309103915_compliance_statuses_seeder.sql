-- +goose Up
-- +goose StatementBegin
INSERT INTO compliance_statuses (title_en, title_ar, description_en, description_ar, color_code, type, status, base_score) VALUES
('Compliant', 'متوافق', 'Compliance document fully complies with regulatory requirements. Regular audits and reviews recommended.', 'مستند الامتثال يتماشى تماماً مع المتطلبات التنظيمية. يوصى بإجراء فحوصات دورية ومراجعات.', '#28A745', 'compliant', 'active', 100),
('Partially Compliant', 'متوافق جزئياً', 'Compliance document is partially compliant. Further actions needed to achieve full compliance.', 'مستند الامتثال جزئياً متوافق. تحتاج إلى إجراءات إضافية لتحقيق الامتثال الكامل.', '#FFC107', 'partially_compliant', 'active', 50),
('Not Implemented', 'لم تنفذ', 'Compliance document is not yet implemented. Action plan required to address compliance gaps.', 'لم يتم تنفيذ مستند الامتثال بعد. تطلب خطة عمل لتغطية الفجوات في الامتثال.', '#DC3545', 'not_implemented', 'active', 0),
('Not Applicable', 'غير قابل للتطبيق', 'Regulatory requirement does not apply to the compliance document. Clear rationale needed.', 'المتطلبات التنظيمية لا تنطبق على مستند الامتثال. الاساس المنطقي الواضح مطلوب.', '#6C757D', 'not_applicable', 'active', 0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE compliance_statuses;
-- +goose StatementEnd
