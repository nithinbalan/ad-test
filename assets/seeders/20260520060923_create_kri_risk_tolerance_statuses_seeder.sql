-- +goose Up
-- +goose StatementBegin
INSERT INTO kri_risk_tolerance_statuses (title_en, title_ar, color_code, type, status) VALUES
    ('Within Minimum Susceptibility', 'ضمن الحد الأدنى من القابلية', '#28A745', 'within_minimum_susceptibility', 'active'),
    ('Within Maximum Susceptibility', 'ضمن الحد الأقصى من القابلية', '#FFC107', 'within_maximum_susceptibility', 'active'),
    ('At Maximum Tolerable', 'عند الحد الأقصى المقبول', '#FD7E14', 'at_maximum_tolerable', 'active'),
    ('Exceeded Maximum Tolerable', 'تجاوز الحد الأقصى المقبول', '#DC3545', 'exceeded_maximum_tolerable', 'active');

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM kri_risk_tolerance_statuses
WHERE type IN (
    'within_minimum_susceptibility',
    'within_maximum_susceptibility',
    'at_maximum_tolerable',
    'exceeded_maximum_tolerable'
);
-- +goose StatementEnd

