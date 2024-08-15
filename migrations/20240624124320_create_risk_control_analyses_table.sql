-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_control_analyses (
    id INT IDENTITY(1,1) PRIMARY KEY,
    type VARCHAR(255) NOT NULL CHECK (type IN ('inherent', 'residual')),
    risk_control_id INT NOT NULL,
    control_test_id INT NULL,
    control_efficiency_status_id INT NOT NULL,
    control_implementation_status_id INT NOT NULL,
    score INT DEFAULT 0,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_risk_control_analyses_risk_controls FOREIGN KEY (risk_control_id) REFERENCES risk_controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_control_analyses_control_tests FOREIGN KEY (control_test_id) REFERENCES control_tests(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_control_analyses_control_efficiency_statuses FOREIGN KEY (control_efficiency_status_id) REFERENCES control_efficiency_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_control_analyses_control_implementation_statuses FOREIGN KEY (control_implementation_status_id) REFERENCES control_implementation_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_control_analyses_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risk_control_analyses;
-- +goose StatementEnd
