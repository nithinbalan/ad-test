-- +goose Up
-- +goose StatementBegin
CREATE TABLE workflow_contexts (
    id INT PRIMARY KEY IDENTITY(1, 1),
    module_id INT NOT NULL,
    type NVARCHAR(255),
    title_en NVARCHAR(MAX) NOT NULL,
    title_ar NVARCHAR(MAX) NOT NULL,
    CONSTRAINT FK_workflow_contexts_modules FOREIGN KEY (module_id) REFERENCES modules (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.workflow_contexts;
-- +goose StatementEnd