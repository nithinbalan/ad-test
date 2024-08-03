-- +goose Up
-- +goose StatementBegin
INSERT INTO risk_timeline_types
    (title_en, title_ar, color_code, type)
VALUES
    ('Identified', 'تم تحديده', '#28A745', 'identified'),
    ('RCA', 'تحليل السبب الجذري', '#007BFF', 'rca'),
    ('Risk Analysis', 'تحليل المخاطر', '#FFC107', 'risk_analysis'),
    ('Control Analysis', 'تحليل التحكم', '#FFA500', 'control_analysis'),
    ('Risk Response', 'الاستجابة للمخاطر', '#6C757D', 'risk_response'),
    ('Residual Control Analysis', 'تحليل التحكم المتبقي', '#D6A2E8', 'residual_control_analysis'),
    ('Residual Risk Analysis', 'تحليل المخاطر المتبقية', '#6F42C1', 'residual_risk_analysis'),
    ('Monitoring', 'مراقبة', '#FD7E14', 'monitoring'),
    ('Closure', 'إغلاق', '#DC3545', 'closure');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_timeline_types;
-- +goose StatementEnd
