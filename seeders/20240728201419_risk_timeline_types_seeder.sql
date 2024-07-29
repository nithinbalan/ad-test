-- +goose Up
-- +goose StatementBegin
INSERT INTO risk_timeline_types
    (title_en, title_ar, color_code, type)
VALUES
    ('Identified', 'تم تحديده', '#FF0000', 'identified'),
    ('RCA', 'تحليل السبب الجذري', '#0000FF', 'rca'),
    ('Risk Analysis', 'تحليل المخاطر', '#FFFF00', 'risk_analysis'),
    ('Control Analysis', 'تحليل التحكم', '#800080', 'control_analysis'),
    ('Risk Response', 'الاستجابة للمخاطر', '#FFA500', 'risk_response'),
    ('Residual Control Analysis', 'تحليل التحكم المتبقي', '#008000', 'residual_control_analysis'),
    ('Residual Risk Analysis', 'تحليل المخاطر المتبقية', '#00FFFF', 'residual_risk_analysis'),
    ('Monitoring', 'مراقبة', '#FF00FF', 'monitoring'),
    ('Closure', 'إغلاق', '#000000', 'closure');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_timeline_types;
-- +goose StatementEnd
