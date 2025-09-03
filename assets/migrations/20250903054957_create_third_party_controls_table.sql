-- +goose Up
-- +goose StatementBegin
CREATE TABLE third_party_controls (
    id INT PRIMARY KEY IDENTITY(1, 1),
	third_party_register_id INT NOT NULL,
    control_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
	CONSTRAINT FK_third_party_controls_third_party_register FOREIGN KEY (third_party_register_id) REFERENCES third_party_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_controls_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_third_party_controls_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE third_party_controls;

-- +goose StatementEnd
