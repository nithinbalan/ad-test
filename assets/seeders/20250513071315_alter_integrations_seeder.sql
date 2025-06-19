-- +goose Up
-- +goose StatementBegin

-- 1. Delete LDAP and Azure AD
DELETE FROM integrations WHERE title_en IN ('LDAP', 'Azure AD','Saudi National ID');

-- 2. Insert new Active Directory
INSERT INTO integrations (
    integration_category_id, title_en, title_ar, description_en, description_ar,
    logo_url, api_base_url, is_enable, created_at, created_by_id, type, status
) VALUES (
    1, 'Active Directory', 'دليل نشط', 
    'Active Directory allows secure connection and authentication to on-premises directory services for managing users and devices.', 
    'يتيح Active Directory الاتصال الآمن بخدمات الدليل المحلية لإدارة المستخدمين والأجهزة.', 
    'Active Directory', NULL, 0, GETDATE(), 1, 'active-directory', 'active'
);

-- 3. Update type for existing integrations based on title_en
UPDATE integrations SET type = 'sharepoint' WHERE title_en = 'SharePoint';
UPDATE integrations SET type = 'onedrive' WHERE title_en = 'OneDrive';
UPDATE integrations SET type = 'google-drive' WHERE title_en = 'Google Drive';
UPDATE integrations SET type = 'tenable' WHERE title_en = 'Tenable';
UPDATE integrations SET type = 'qualys' WHERE title_en = 'Qualys';
UPDATE integrations SET type = 'manageengine' WHERE title_en = 'ManageEngine';
UPDATE integrations SET type = 'ibm' WHERE title_en = 'IBM';
UPDATE integrations SET type = 'oracle' WHERE title_en = 'Oracle';
UPDATE integrations SET type = 'odoo' WHERE title_en = 'Odoo';
UPDATE integrations SET type = 'jira' WHERE title_en = 'Jira';
UPDATE integrations SET type = 'chatgpt' WHERE title_en = 'ChatGPT';
UPDATE integrations SET type = 'gemini' WHERE title_en = 'Gemini';
UPDATE integrations SET type = 'notion' WHERE title_en = 'Notion';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- 1. Delete newly inserted Active Directory
DELETE FROM integrations WHERE title_en = 'Active Directory';

-- 2. (Optional) Re-insert LDAP and Azure AD
INSERT INTO integrations (
    integration_category_id, title_en, title_ar, description_en, description_ar,
    logo_url, api_base_url, is_enable, created_at, created_by_id
) VALUES
(1,'LDAP', 'لداب', 'LDAP is a protocol for managing directory services over a network.', 'LDAP هو بروتوكول لإدارة خدمات الدليل عبر الشبكة.', 'LDAP', NULL, 0, GETDATE(), 1),
(1,'Azure AD', 'أزور م', 'Active Directory not built for remote work-try JumpCloud free.', 'Active Directory لم يتم تصميمه للعمل عن بعد - جرب JumpCloud مجانًا.', 'Azure AD', NULL, 0, GETDATE(), 1);

-- 3. Optional: Clear type column values (if you want to reverse)
UPDATE integrations SET type = NULL WHERE title_en IN (
    'Saudi National ID',
    'SharePoint',
    'OneDrive',
    'Google Drive',
    'Tenable',
    'Qualys',
    'ManageEngine',
    'IBM',
    'Oracle',
    'Odoo',
    'Jira',
    'ChatGPT',
    'Gemini',
    'Notion'
);

-- +goose StatementEnd
