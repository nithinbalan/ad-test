-- +goose Up
-- +goose StatementBegin
INSERT INTO compliance_assessment_impact_types
    (title_en, title_ar, type, description_en, description_ar, status)
VALUES
    (
        'Regulatory',
        'الأثر التنظيمي',
        'regulatory',
        'Impact related to non-compliance with laws, regulations, or statutory requirements.',
        'الأثر الناتج عن عدم الامتثال للقوانين أو اللوائح أو المتطلبات النظامية.',
        'active'
    ),
    (
        'Reputational',
        'الأثر على السمعة',
        'reputational',
        'Impact affecting the organization’s public image, trust, or credibility.',
        'الأثر الذي يؤثر على صورة المؤسسة العامة أو الثقة أو المصداقية.',
        'active'
    ),
    (
        'Financial',
        'الأثر المالي',
        'financial',
        'Impact resulting in financial loss, penalties, fines, or increased operational costs.',
        'الأثر الناتج عن خسائر مالية أو غرامات أو تكاليف تشغيلية إضافية.',
        'active'
    ),
    (
        'Operational',
        'الأثر التشغيلي',
        'operational',
        'Impact disrupting business operations, processes, systems, or service delivery.',
        'الأثر الذي يعطل العمليات التجارية أو الأنظمة أو تقديم الخدمات.',
        'active'
    );
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM compliance_assessment_impact_types
WHERE type IN ('regulatory', 'reputational', 'financial', 'operational');
-- +goose StatementEnd