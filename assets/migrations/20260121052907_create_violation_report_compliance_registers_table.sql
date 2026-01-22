-- +goose Up
-- +goose StatementBegin
CREATE TABLE violation_report_compliance_registers
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    violation_report_id INT NOT NULL,
    compliance_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_violation_report_compliance_registers_violation_report FOREIGN KEY (violation_report_id) REFERENCES violation_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_report_compliance_registers_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE violation_report_compliance_registers;

-- +goose StatementEnd