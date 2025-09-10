-- +goose Up
-- +goose StatementBegin

ALTER TABLE workflows
DROP CONSTRAINT CHK_workflows_type;

ALTER TABLE workflows
ALTER COLUMN type VARCHAR(100) NULL;

ALTER TABLE workflows
ADD workflow_context_id INT NULL;

ALTER TABLE workflows
ADD CONSTRAINT FK_workflows_workflow_context
    FOREIGN KEY (workflow_context_id) REFERENCES workflow_contexts(id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE workflows
DROP CONSTRAINT FK_workflows_workflow_context;

ALTER TABLE workflows
DROP COLUMN workflow_context_id;

ALTER TABLE workflows
ALTER COLUMN type VARCHAR(100) NOT NULL;

ALTER TABLE workflows
ADD CONSTRAINT CHK_workflows_type CHECK (type IN ('document', 'draft-report', 'final-report'));

-- +goose StatementEnd
