-- +goose Up
-- +goose StatementBegin

ALTER TABLE policy_exceptions
ADD document_requirement_id INT NULL;

ALTER TABLE policy_exceptions
ADD CONSTRAINT FK_policy_exceptions_document_requirements
FOREIGN KEY (document_requirement_id) REFERENCES document_requirements(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE policy_exceptions
DROP CONSTRAINT FK_policy_exceptions_document_requirements;

ALTER TABLE policy_exceptions
DROP COLUMN document_requirement_id;

-- +goose StatementEnd
