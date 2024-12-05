-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_threats
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    risk_register_id INT NOT NULL,
    threat_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_risk_threats_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_threats_threat FOREIGN KEY (threat_id) REFERENCES threats(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_threats_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_threats;

-- +goose StatementEnd
