-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_milestone_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description_en NVARCHAR(MAX),
    description_ar NVARCHAR(MAX),
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'pending',
        CONSTRAINT CHK_project_milestone_statuses_type CHECK (type IN ('pending', 'in_progress', 'completed', 'cancelled')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_project_milestone_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.project_milestone_statuses;
-- +goose StatementEnd
