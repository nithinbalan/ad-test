-- +goose Up
-- +goose StatementBegin
INSERT INTO compliance_assessment_framework_types (
    type, title_en, title_ar, description_en, description_ar, status
) VALUES
('binary', 'Binary Framework', 'إطار ثنائي', 'Basic Yes/No scoring', 'تسجيل نعم/لا الأساسي', 'active'),
('maturity', 'Maturity Framework', 'إطار النضج', 'Uses maturity stages', 'يستخدم مراحل النضج', 'active'),
('point', 'Point Framework', 'إطار النقاط', 'Uses point allocation', 'يستخدم تخصيص النقاط', 'active'),
('percentage', 'Percentage Framework', 'إطار النسبة المئوية', 'Uses implementation %', 'يستخدم نسبة التنفيذ المئوية', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE compliance_assessment_framework_types;
-- +goose StatementEnd
