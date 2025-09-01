-- +goose Up
-- +goose StatementBegin
CREATE TABLE  third_party_risks (
    id INT PRIMARY KEY IDENTITY(1, 1),
	third_party_register_id INT NOT NULL,
    risk_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
	CONSTRAINT FK_third_party_risks_third_party_register FOREIGN KEY (third_party_register_id) REFERENCES third_party_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_risks_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_third_party_risks_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE  third_party_risks;

-- +goose StatementEnd
