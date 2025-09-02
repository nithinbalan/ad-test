
-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents
ADD policy_folder_id INT NULL;

ALTER TABLE documents
ADD CONSTRAINT FK_documents_policy_folders
FOREIGN KEY (policy_folder_id) REFERENCES policy_folders(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE documents
DROP CONSTRAINT FK_documents_policy_folders;


ALTER TABLE documents
DROP COLUMN policy_folder_id;
-- +goose StatementEnd
