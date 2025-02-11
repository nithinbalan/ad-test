-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_register
ADD compliance_folder_id INT NULL;

ALTER TABLE compliance_register
ADD CONSTRAINT FK_compliance_register_compliance_folder FOREIGN KEY (compliance_folder_id) REFERENCES compliance_folders(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_register
DROP CONSTRAINT FK_compliance_register_compliance_folder;

ALTER TABLE compliance_register
DROP COLUMN compliance_folder_id;
-- +goose StatementEnd
