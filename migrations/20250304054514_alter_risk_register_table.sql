-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_register
ADD submitted_by INT NULL,
    submitted_at DATETIME NULL,
    reviewed_by INT NULL,
    CONSTRAINT FK_risk_register_users_4 FOREIGN KEY (submitted_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_users_5 FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_register
DROP CONSTRAINT FK_risk_register_users_4;

ALTER TABLE risk_register
DROP CONSTRAINT FK_risk_register_users_5;

ALTER TABLE risk_register
DROP COLUMN submitted_by,
DROP COLUMN submitted_at,
DROP COLUMN reviewed_by;
-- +goose StatementEnd
