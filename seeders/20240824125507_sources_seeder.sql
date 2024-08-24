-- +goose Up
-- +goose StatementBegin
INSERT INTO sources
    (title_en, title_ar, type, status)
VALUES
    ('Direct', 'مباشر', 'direct', 'active'),
    ('Jethur Document Management', 'إدارة مستندات Jethur', 'jethur-dm', 'active'),
    ('Google Drive', 'جوجل درايف', 'gdrive', 'active'),
    ('OneDrive', 'ون درايف', 'onedrive', 'active'),
    ('SharePoint', 'شيربوينت', 'sharepoint', 'active'),
    ('Dropbox', 'دروب بوكس', 'dropbox', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE sources;
-- +goose StatementEnd
