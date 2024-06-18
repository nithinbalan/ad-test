-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_controls (
    id INT PRIMARY KEY IDENTITY(1, 1),
    risk_register_id INT NOT NULL,
    control_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_risk_controls_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_controls_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_controls;

-- +goose StatementEnd