-- +goose Up
-- +goose StatementBegin
ALTER TABLE processes ADD section_id INT NULL;
ALTER TABLE processes ADD sub_section_id INT NULL;

ALTER TABLE processes
ADD CONSTRAINT FK_processes_section
FOREIGN KEY (section_id) REFERENCES organization_units(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE processes
ADD CONSTRAINT FK_processes_sub_section
FOREIGN KEY (sub_section_id) REFERENCES organization_units(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE processes DROP CONSTRAINT FK_processes_sub_section;
ALTER TABLE processes DROP CONSTRAINT FK_processes_section;
ALTER TABLE processes DROP COLUMN sub_section_id;
ALTER TABLE processes DROP COLUMN section_id;
-- +goose StatementEnd
