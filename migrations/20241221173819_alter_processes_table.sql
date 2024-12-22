-- +goose Up
-- +goose StatementBegin
ALTER TABLE processes 
ADD is_bia_performed TINYINT NOT NULL DEFAULT 0;

ALTER TABLE processes 
ADD last_assessed_at DATETIME NULL;

ALTER TABLE processes 
ADD last_assessed_by_id INT NULL;

ALTER TABLE processes
ADD CONSTRAINT FK_processes_users_4 FOREIGN KEY (last_assessed_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE processes DROP CONSTRAINT FK_processes_users_4;
ALTER TABLE processes DROP COLUMN last_assessed_by_id;
ALTER TABLE processes DROP COLUMN last_assessed_at;
ALTER TABLE processes DROP COLUMN is_bia_performed;
-- +goose StatementEnd
