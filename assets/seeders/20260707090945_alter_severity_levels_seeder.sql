-- +goose Up
-- +goose StatementBegin
UPDATE severity_levels
SET weight = 4
WHERE type = 'high';

UPDATE severity_levels
SET weight = 3
WHERE type = 'medium';

UPDATE severity_levels
SET weight = 2
WHERE type = 'low';

IF NOT EXISTS (
    SELECT 1
    FROM severity_levels
    WHERE type = 'critical'
)
BEGIN
    INSERT INTO severity_levels (
        title_en,
        title_ar,
        color_code,
        type,
        status,
        weight
    )
    VALUES (
        'Critical',
        N'حرج',
        '#E0554D',
        'critical',
        'active',
        5
    );
END

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM severity_levels
WHERE type = 'critical';

UPDATE severity_levels
SET weight = 10
WHERE type = 'high';

UPDATE severity_levels
SET weight = 5
WHERE type = 'medium';

UPDATE severity_levels
SET weight = 2
WHERE type = 'low';
-- +goose StatementEnd