-- +goose Up
-- +goose StatementBegin
ALTER TABLE audit_document_requests
ADD delegatee_user_id INT NULL;

ALTER TABLE audit_document_requests
ADD CONSTRAINT FK_audit_document_requests_delegatee_user
FOREIGN KEY (delegatee_user_id) REFERENCES users(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE audit_document_requests
DROP CONSTRAINT FK_audit_document_requests_delegatee_user;

ALTER TABLE audit_document_requests
DROP COLUMN delegatee_user_id;
-- +goose StatementEnd
