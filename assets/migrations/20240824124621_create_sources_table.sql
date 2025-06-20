-- +goose Up
-- +goose StatementBegin
CREATE TABLE sources
(
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'direct',
    CONSTRAINT CHK_sources_type CHECK (type IN ('direct', 'jethur-dm', 'gdrive', 'onedrive', 'sharepoint', 'dropbox')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_sources_status CHECK (status IN ('active', 'inactive'))

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sources;
-- +goose StatementEnd
