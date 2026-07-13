-- +goose Up
-- +goose StatementBegin
ALTER TABLE sla_escalation_matrix
ADD notifiying_role_id INT NULL;

ALTER TABLE sla_escalation_matrix ADD CONSTRAINT FK_sla_escalation_matrix_notifiying_role
FOREIGN KEY (notifiying_role_id) REFERENCES roles(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sla_escalation_matrix
DROP CONSTRAINT FK_sla_escalation_matrix_notifiying_role;

ALTER TABLE sla_escalation_matrix DROP COLUMN notifiying_role_id;
-- +goose StatementEnd