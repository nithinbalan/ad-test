-- +goose Up
-- +goose StatementBegin
UPDATE impact_categories
SET type = 'financial'
WHERE title_en = 'Financial';

UPDATE impact_categories
SET type = 'operational'
WHERE title_en = 'Operational';

UPDATE impact_categories
SET type = 'reputational'
WHERE title_en = 'Reputational';

UPDATE impact_categories
SET type = 'legal'
WHERE title_en = 'Legal/Compliance';

UPDATE impact_categories
SET type = 'human-resource'
WHERE title_en = 'Human Resources';

UPDATE impact_categories
SET type = 'it'
WHERE title_en = 'IT/System';

UPDATE impact_categories
SET type = 'supply-chain'
WHERE title_en = 'Supply Chain';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert the type back to its previous value 
UPDATE impact_categories
SET type = NULL
WHERE type = 'Financial';

UPDATE impact_categories
SET type = NULL
WHERE type = 'Operational';

UPDATE impact_categories
SET type = NULL
WHERE type = 'Reputational';

UPDATE impact_categories
SET type = NULL
WHERE type = 'Legal/Compliance';

UPDATE impact_categories
SET type = NULL
WHERE type = 'Human Resources';

UPDATE impact_categories
SET type = NULL
WHERE type = 'IT/System';

UPDATE impact_categories
SET type = NULL
WHERE type = 'Supply Chain';
-- +goose StatementEnd
