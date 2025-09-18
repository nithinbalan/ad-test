-- +goose Up
-- +goose StatementBegin
ALTER TABLE activity_logs
ADD 
    module_id INT NULL,
    sub_module_id INT,
    section NVARCHAR(255);

ALTER TABLE activity_logs
ADD CONSTRAINT fk_activity_logs_module
FOREIGN KEY (module_id) REFERENCES modules(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

ALTER TABLE activity_logs
ADD CONSTRAINT fk_activity_logs_sub_module
FOREIGN KEY (sub_module_id) REFERENCES modules(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE activity_logs DROP CONSTRAINT fk_activity_logs_sub_module;
ALTER TABLE activity_logs DROP CONSTRAINT fk_activity_logs_module;

ALTER TABLE activity_logs
DROP COLUMN module_id,
DROP COLUMN sub_module_id,
DROP COLUMN section;
-- +goose StatementEnd
