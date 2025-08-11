-- +goose Up
-- +goose StatementBegin
CREATE TABLE third_party_strategic_objectives  (
    id INT PRIMARY KEY IDENTITY(1,1),
    strategic_objective_id INT NOT NULL,
    third_party_register_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    CONSTRAINT FK_third_party_strategic_objectives_strategic_objectives FOREIGN KEY (strategic_objective_id) REFERENCES strategic_objectives(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_strategic_objectives_third_party_register FOREIGN KEY (third_party_id) REFERENCES third_party_register(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_strategic_objectives_user_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE third_party_strategic_objectives;
-- +goose StatementEnd
