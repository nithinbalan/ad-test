-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks ADD document_change_request_id INT NULL;
ALTER TABLE tasks ADD document_change_request_workflow_id INT NULL;

ALTER TABLE tasks ADD CONSTRAINT FK_tasks_document_change_request 
    FOREIGN KEY (document_change_request_id) REFERENCES document_change_requests(id) ON DELETE SET NULL ON UPDATE NO ACTION;

ALTER TABLE tasks ADD CONSTRAINT FK_tasks_document_change_request_workflow 
    FOREIGN KEY (document_change_request_workflow_id) REFERENCES document_change_request_workflows(id) ON DELETE SET NULL ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks DROP CONSTRAINT FK_tasks_document_change_request;
ALTER TABLE tasks DROP CONSTRAINT FK_tasks_document_change_request_workflow;
ALTER TABLE tasks DROP COLUMN document_change_request_id;
ALTER TABLE tasks DROP COLUMN document_change_request_workflow_id;
-- +goose StatementEnd
