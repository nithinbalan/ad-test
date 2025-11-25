-- +goose Up
-- +goose StatementBegin
ALTER TABLE integrations
    DROP CONSTRAINT IF EXISTS CHK_integrations_type;

ALTER TABLE integrations
    ADD CONSTRAINT CHK_integrations_type
        CHECK (type IN (
            'active-directory',
            'azure-ad',
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
ALTER TABLE integrations
    DROP CONSTRAINT IF EXISTS CHK_integrations_type;

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
