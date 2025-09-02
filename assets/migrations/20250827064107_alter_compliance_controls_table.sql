-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_controls
ADD compliance_status_id INT NULL;

ALTER TABLE compliance_controls
ADD CONSTRAINT FK_compliance_controls_compliance_status
    FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses(id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_controls
DROP CONSTRAINT FK_compliance_controls_compliance_status;

ALTER TABLE compliance_controls
DROP COLUMN compliance_status_id;
-- +goose StatementEnd
