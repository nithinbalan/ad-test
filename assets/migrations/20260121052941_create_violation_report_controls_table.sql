-- +goose Up
-- +goose StatementBegin
CREATE TABLE violation_report_controls (
    id INT PRIMARY KEY IDENTITY(1, 1),
    violation_report_id INT NOT NULL,
    control_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_violation_report_controls_violation_report FOREIGN KEY (violation_report_id) REFERENCES violation_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_report_controls_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE violation_report_controls;

-- +goose StatementEnd