-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description_en NVARCHAR(MAX),
    description_ar NVARCHAR(MAX),
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'not_started',
        CONSTRAINT CHK_project_statuses_type CHECK (type IN ('not_started','planning','in_progress','on_hold','cancelled','completed')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_project_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.project_statuses;
-- +goose StatementEnd