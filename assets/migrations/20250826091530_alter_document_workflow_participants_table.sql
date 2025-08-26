-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_workflow_participants
ADD delegatee_user_id INT NULL;

ALTER TABLE document_workflow_participants
ADD CONSTRAINT FK_document_workflow_participants_delegatee_user
FOREIGN KEY (delegatee_user_id) REFERENCES users(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_workflow_participants
DROP CONSTRAINT FK_document_workflow_participants_delegatee_user;

ALTER TABLE document_workflow_participants
DROP COLUMN delegatee_user_id;
-- +goose StatementEnd
