-- +goose Up
-- +goose StatementBegin

-- Step 1: Add the type column
ALTER TABLE integrations
ADD type NVARCHAR(100) NULL;

-- +goose StatementEnd

-- +goose StatementBegin

-- Step 2: Add the status column (default 'inactive')
ALTER TABLE integrations
ADD status VARCHAR(10) NOT NULL DEFAULT 'inactive';

-- +goose StatementEnd

-- +goose StatementBegin

-- Step 3: Add the check constraint on type
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

-- +goose StatementBegin

-- Step 4: Add the check constraint on status
ALTER TABLE integrations
ADD CONSTRAINT CHK_integrations_status 
CHECK (status IN ('active', 'inactive'));

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Step 1: Drop constraint on status
ALTER TABLE integrations
DROP CONSTRAINT CHK_integrations_status;

-- +goose StatementEnd

-- +goose StatementBegin

-- Step 2: Drop constraint on type
ALTER TABLE integrations
DROP CONSTRAINT CHK_integrations_type;

-- +goose StatementEnd

-- +goose StatementBegin

-- Step 3: Drop status column
ALTER TABLE integrations
DROP COLUMN status;

-- +goose StatementEnd

-- +goose StatementBegin

-- Step 4: Drop type column
ALTER TABLE integrations
DROP COLUMN type;

-- +goose StatementEnd
