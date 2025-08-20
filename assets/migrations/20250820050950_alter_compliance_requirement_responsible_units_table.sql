-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_requirement_responsible_units
ADD completed_at DATETIME NULL,
    completed_by_user_id INT NULL,
    comments NVARCHAR(MAX) NULL;

ALTER TABLE compliance_requirement_responsible_units
ADD CONSTRAINT FK_crru_completed_by FOREIGN KEY (completed_by_user_id)
    REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_requirement_responsible_units
DROP CONSTRAINT FK_crru_completed_by;

ALTER TABLE compliance_requirement_responsible_units
DROP COLUMN completed_at,
    completed_by_user_id,
    comments;
-- +goose StatementEnd
