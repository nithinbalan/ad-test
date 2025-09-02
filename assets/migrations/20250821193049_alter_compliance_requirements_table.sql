-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_requirements
DROP CONSTRAINT FK_compliance_requirements_users_1;

ALTER TABLE compliance_requirements
DROP COLUMN responsible_user_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_requirements
ADD responsible_user_id INT NULL;

ALTER TABLE compliance_requirements
ADD CONSTRAINT FK_compliance_requirements_users_1
FOREIGN KEY (responsible_user_id) REFERENCES users(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd
