-- +goose Up
-- +goose StatementBegin
CREATE TABLE workflow_level_types (
    id INT PRIMARY KEY IDENTITY(1,1),
    type VARCHAR(20) NOT NULL,
     CONSTRAINT CHK_workflow_level_types_type CHECK (type IN ('individual', 'multiple-users', 'role', 'hod')),
    title_en VARCHAR(45) NOT NULL,
    title_ar VARCHAR(45) NOT NULL,
    status VARCHAR(20) NOT NULL,
     CONSTRAINT CHK_workflow_level_types_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.workflow_level_types;
-- +goose StatementEnd
