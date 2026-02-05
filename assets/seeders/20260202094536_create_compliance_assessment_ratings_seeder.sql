-- +goose Up
-- +goose StatementBegin
INSERT INTO compliance_assessment_ratings
    (title_en, title_ar, type, color_code, description_en, description_ar, score_from, score_to, status)
VALUES
    (
        'Critical',
        'حرج',
        'critical',
        '#B22222',
        NULL,
        NULL,
        3.00,
        4.00,
        'active'
    ),
    (
        'High',
        'مرتفع',
        'high',
        '#FF0000',
        NULL,
        NULL,
        2.00,
        2.99,
        'active'
    ),
    (
        'Medium',
        'متوسط',
        'medium',
        '#FFA500',
        NULL,
        NULL,
        1.00,
        1.99,
        'active'
    ),
    (
        'Low',
        'منخفض',
        'low',
        '#90DD75',
        NULL,
        NULL,
        0.01,
        0.99,
        'active'
    ),
    (
        'Closed',
        'مغلق',
        'closed',
        '#008000',
        NULL,
        NULL,
        0.00,
        0.00,
        'active'
    );
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM compliance_assessment_ratings
WHERE type IN ('critical', 'high', 'medium', 'low', 'closed');
-- +goose StatementEnd