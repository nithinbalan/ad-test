-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_register
ADD severity_level_id INT NULL;

ALTER TABLE risk_register
ADD CONSTRAINT FK_risk_register_severity_levels
FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_register
DROP CONSTRAINT FK_risk_register_severity_levels;


ALTER TABLE risk_register
DROP COLUMN severity_level_id;
-- +goose StatementEnd
