-- +goose Up
-- +goose StatementBegin
UPDATE integrations SET logo_url = 'active-directory' WHERE type = 'active-directory';
UPDATE integrations SET logo_url = 'sharepoint' WHERE type = 'sharepoint';
UPDATE integrations SET logo_url = 'one-drive' WHERE type = 'onedrive';
UPDATE integrations SET logo_url = 'google-drive' WHERE type = 'google-drive';
UPDATE integrations SET logo_url = 'tenable' WHERE type = 'tenable';
UPDATE integrations SET logo_url = 'qualys' WHERE type = 'qualys';
UPDATE integrations SET logo_url = 'manageengine' WHERE type = 'manageengine';
UPDATE integrations SET logo_url = 'ibm' WHERE type = 'ibm';
UPDATE integrations SET logo_url = 'oracle' WHERE type = 'oracle';
UPDATE integrations SET logo_url = 'odoo' WHERE type = 'odoo';
UPDATE integrations SET logo_url = 'jira' WHERE type = 'jira';
UPDATE integrations SET logo_url = 'chatgpt' WHERE type = 'chatgpt';
UPDATE integrations SET logo_url = 'gemini' WHERE type = 'gemini';
UPDATE integrations SET logo_url = 'notion' WHERE type = 'notion';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Rollback: Reset logo_url to NULL
UPDATE integrations SET logo_url = NULL;
-- +goose StatementEnd
