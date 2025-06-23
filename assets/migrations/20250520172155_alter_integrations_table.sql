-- +goose Up
-- +goose StatementBegin

-- Step: Drop the existing check constraint on type
ALTER TABLE integrations
DROP CONSTRAINT CHK_integrations_type;

-- +goose StatementEnd

-- +goose StatementBegin

-- Step: Re-add the check constraint including 'siem'
ALTER TABLE integrations
ADD CONSTRAINT CHK_integrations_type 
CHECK (type IN (
    'active-directory',
    'sharepoint',
    'onedrive',
    'google-drive',
    'tenable',
    'qualys',
    'manageengine',
    'ibm',
    'oracle',
    'odoo',
    'jira',
    'chatgpt',
    'gemini',
    'notion',
    'siem'
));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Step: Drop the updated constraint with 'siem'
ALTER TABLE integrations
DROP CONSTRAINT CHK_integrations_type;

-- +goose StatementEnd

-- +goose StatementBegin

-- Step: Re-add the original constraint (without 'siem')
ALTER TABLE integrations
ADD CONSTRAINT CHK_integrations_type 
CHECK (type IN (
    'active-directory',
    'sharepoint',
    'onedrive',
    'google-drive',
    'tenable',
    'qualys',
    'manageengine',
    'ibm',
    'oracle',
    'odoo',
    'jira',
    'chatgpt',
    'gemini',
    'notion'
));

-- +goose StatementEnd
