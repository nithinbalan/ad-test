-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_report_risk_registers (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_report_id INT NOT NULL,
    risk_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_risk_report_risk_registers_risk_reports FOREIGN KEY (risk_report_id) REFERENCES risk_reports (id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_report_risk_registers_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_report_risk_registers_users FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_report_risk_registers;
-- +goose StatementEnd
