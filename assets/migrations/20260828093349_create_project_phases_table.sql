-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_phases (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
        CONSTRAINT CHK_project_phases_type CHECK (type IN ('initiation', 'planning', 'execution', 'go_live', 'closure')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_project_phases_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.project_phases;
-- +goose StatementEnd
