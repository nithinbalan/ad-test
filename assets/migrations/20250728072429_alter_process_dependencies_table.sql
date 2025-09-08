-- +goose Up
-- +goose StatementBegin
ALTER TABLE process_dependencies
ADD third_party_register_id INT NULL;

ALTER TABLE process_dependencies
ADD CONSTRAINT FK_process_dependencies_third_party_register
FOREIGN KEY (third_party_register_id)
REFERENCES third_party_register(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE process_dependencies
DROP CONSTRAINT FK_process_dependencies_third_party_register;

ALTER TABLE process_dependencies
DROP COLUMN third_party_register_id;
-- +goose StatementEnd
