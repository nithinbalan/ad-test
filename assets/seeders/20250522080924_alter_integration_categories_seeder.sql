-- +goose Up
-- +goose StatementBegin

-- Update integration category titles
UPDATE integration_categories
SET title_en = 'Single Sign-On (SSO)'
WHERE title_en = 'Single sign-on';

UPDATE integration_categories
SET title_en = 'Security Information and Event Management (SIEM)'
WHERE title_en = 'SIEM';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert integration category titles
UPDATE integration_categories
SET title_en = 'Single sign-on'
WHERE title_en = 'Single Sign-On (SSO)';

UPDATE integration_categories
SET title_en = 'SIEM'
WHERE title_en = 'Security Information and Event Management (SIEM)';

-- +goose StatementEnd
