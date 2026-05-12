-- +goose Up
-- +goose StatementBegin

CREATE TABLE risk_update_plan_timeline_types
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(500) NOT NULL,
    CONSTRAINT CHK_risk_update_plan_timeline_types_type CHECK (type IN (
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
    ))
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_update_plan_timeline_types;
-- +goose StatementEnd

