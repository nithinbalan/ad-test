-- +goose Up
-- +goose StatementBegin

DELETE FROM risk_response_strategies 
WHERE type = 'tolerate';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

INSERT INTO risk_response_strategies (type, title_en, title_ar, description_en, description_ar, status)
VALUES (
    'tolerate', 
    'Tolerate', 
    N'تحمل', 
    'Risk event that is accepted but monitored for changes in conditions.', 
    N'حدث المخاطر الذي يتم قبوله ولكن يتم مراقبته للتغيرات في الظروف.', 
    'active'
);

-- +goose StatementEnd