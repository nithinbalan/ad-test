-- +goose Up
-- +goose StatementBegin
CREATE TABLE process_dependency_sub_types (
    id INT PRIMARY KEY IDENTITY(1, 1),
    process_dependency_type_id INT NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_process_dependency_sub_types_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_process_dependency_sub_types_process_dependency_types FOREIGN KEY (process_dependency_type_id) REFERENCES process_dependency_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.process_dependency_sub_types;
-- +goose StatementEnd
