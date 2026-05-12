-- +goose Up
-- +goose StatementBegin

INSERT INTO risk_update_plan_timeline_types (title_en, title_ar, color_code, type)
VALUES
('Risk Register Update Plan', N'خطة تحديث سجل المخاطر', '#90DD75', 'risk_register_update_plan'),
('ERM Department Evaluation', N'تقييم قسم إدارة المخاطر', '#90DD75', 'erm_department_evaluation'),
('Introductory Meeting', N'الاجتماع التمهيدي', '#90DD75', 'introductory_meeting'),
('Business Context', N'سياق العمل', '#90DD75', 'business_context'),
('Self-Awareness Training', N'تدريب الوعي الذاتي', '#90DD75', 'self_awareness_training'),
('Leadership Assessment', N'تقييم القيادة', '#90DD75', 'leadership_assessment'),
('Risk Team Assessment', N'تقييم فريق المخاطر', '#90DD75', 'risk_team_assessment'),
('Risk Identification', N'تحديد المخاطر', '#90DD75', 'risk_identification'),

('Control Register Updated', N'تحديث سجل الضوابط', '#FFFF00', 'control_register_updated'),

('Risk Assessment Workshop', N'ورشة تقييم المخاطر', '#FFC879', 'risk_assessment_workshop'),
('Risk Workshop Assessment', N'تقييم ورشة المخاطر', '#FFC879', 'risk_workshop_assessment'),
('Risk Analysis', N'تحليل المخاطر', '#FFC879', 'risk_analysis'),

('Risk Control Association', N'ارتباط الضوابط', '#0E7B61', 'risk_control_association'),

('Risk Mitigation', N'تخفيف المخاطر', '#DB0015', 'risk_mitigation');

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM risk_update_plan_timeline_types
WHERE type IN (
    'risk_register_update_plan',
    'erm_department_evaluation',
    'introductory_meeting',
    'business_context',
    'self_awareness_training',
    'leadership_assessment',
    'risk_team_assessment',
    'risk_identification',
    'control_register_updated',
    'risk_assessment_workshop',
    'risk_workshop_assessment',
    'risk_analysis',
    'risk_control_association',
    'risk_mitigation'
);

-- +goose StatementEnd