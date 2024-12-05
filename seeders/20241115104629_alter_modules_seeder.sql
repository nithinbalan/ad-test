-- +goose Up
-- +goose StatementBegin
UPDATE modules
SET title_en = 'General', title_ar = 'عام'
WHERE identifier = 'general_setting';

UPDATE modules
SET title_en = 'Organization', title_ar = 'منظمة'
WHERE identifier = 'organization_setting';

UPDATE modules
SET title_en = 'Compliance', title_ar = 'امتثال'
WHERE identifier = 'compliance_setting';

UPDATE modules
SET title_en = 'Risk', title_ar = 'مخاطرة'
WHERE identifier = 'risk_setting';

UPDATE modules
SET title_en = 'Audit', title_ar = 'مراجعة'
WHERE identifier = 'audit_setting';

UPDATE modules
SET title_en = 'Document', title_ar = 'وثيقة'
WHERE identifier = 'document_setting';

UPDATE modules
SET title_en = 'Incident', title_ar = 'حادثة'
WHERE identifier = 'incident_setting';

UPDATE modules
SET title_en = 'Review Meeting', title_ar = 'اجتماع المراجعة'
WHERE identifier = 'review_meeting_setting';

UPDATE modules
SET title_en = 'Audit Insight', title_ar = 'رؤية التدقيق'
WHERE identifier = 'audit_insight_setting';

UPDATE modules
SET title_en = 'Business Continuity', title_ar = 'استمرارية الأعمال'
WHERE identifier = 'business_continuity_setting';

UPDATE modules
SET title_en = 'Asset', title_ar = 'أصل'
WHERE identifier = 'asset_setting';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules
SET title_en = 'General Setting', title_ar = 'الإعداد العام'
WHERE identifier = 'general_setting';

UPDATE modules
SET title_en = 'Organization Setting', title_ar = 'إعداد المنظمة'
WHERE identifier = 'organization_setting';

UPDATE modules
SET title_en = 'Compliance Setting', title_ar = 'إعداد الامتثال'
WHERE identifier = 'compliance_setting';

UPDATE modules
SET title_en = 'Risk Setting', title_ar = 'تحديد المخاطر'
WHERE identifier = 'risk_setting';

UPDATE modules
SET title_en = 'Audit Setting', title_ar = 'إعداد التدقيق'
WHERE identifier = 'audit_setting';

UPDATE modules
SET title_en = 'Document Setting', title_ar = 'إعداد الوثيقة'
WHERE identifier = 'document_setting';

UPDATE modules
SET title_en = 'Incident Setting', title_ar = 'إعداد الحادث'
WHERE identifier = 'incident_setting';

UPDATE modules
SET title_en = 'Review Meeting Setting', title_ar = 'مراجعة إعدادات الاجتماع'
WHERE identifier = 'review_meeting_setting';

UPDATE modules
SET title_en = 'Audit Insight Setting', title_ar = 'إعداد رؤية التدقيق'
WHERE identifier = 'audit_insight_setting';

UPDATE modules
SET title_en = 'Business Continuity Setting', title_ar = 'إعداد استمرارية الأعمال'
WHERE identifier = 'business_continuity_setting';

UPDATE modules
SET title_en = 'Asset Setting', title_ar = 'إعداد الأصول'
WHERE identifier = 'asset_setting';
-- +goose StatementEnd
