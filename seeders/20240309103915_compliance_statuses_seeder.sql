-- +goose Up
-- +goose StatementBegin
INSERT INTO compliance_statuses (title_en, title_ar, description_en, description_ar, color_code, type, status, base_score) VALUES
('Not Implemented', 'لم تنفذ', 'Compliance document is not yet implemented. Action plan required to address compliance gaps.', 'لم يتم تنفيذ مستند الامتثال بعد. تطلب خطة عمل لتغطية الفجوات في الامتثال.', '#DC3545', 'not_implemented', 'active', 0),
('In Progress', 'في تَقَدم', 'Compliance document implementation is underway. Ongoing monitoring and evaluation required.', 'تم بدء تنفيذ مستند الامتثال. الرصد والتقييم المستمرين مطلوبين.', '#FFC107', 'in_progress', 'active', 25),
('Partially Compliant', 'متوافق جزئياً', 'Compliance document is partially compliant. Further actions needed to achieve full compliance.', 'مستند الامتثال جزئياً متوافق. تحتاج إلى إجراءات إضافية لتحقيق الامتثال الكامل.', '#FFC107', 'partially_compliant', 'active', 50),
('Compliant', 'متوافق', 'Compliance document fully complies with regulatory requirements. Regular audits and reviews recommended.', 'مستند الامتثال يتماشى تماماً مع المتطلبات التنظيمية. يوصى بإجراء فحوصات دورية ومراجعات.', '#28A745', 'compliant', 'active', 100),
('Exempt', 'معفى', 'Compliance document exempted from specific regulatory requirements. Justification documentation required.', 'مستند الامتثال معفى من متطلبات التنظيم الخاصة. الوثائق التوضيحية مطلوبة.', '#17A2B8', 'exempt', 'active', 0),
('Not Applicable', 'غير قابل للتطبيق', 'Regulatory requirement does not apply to the compliance document. Clear rationale needed.', 'المتطلبات التنظيمية لا تنطبق على مستند الامتثال. الاساس المنطقي الواضح مطلوب.', '#6C757D', 'not_applicable', 'active', 0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE compliance_statuses;
-- +goose StatementEnd
