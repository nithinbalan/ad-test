-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategic_objective_controls (
    id INT PRIMARY KEY IDENTITY(1, 1),
    control_id INT NOT NULL,
    strategic_objective_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_strategic_objective_control_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategic_objective_control_strategic_objectives FOREIGN KEY (strategic_objective_id) REFERENCES strategic_objectives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategic_objective_control_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE strategic_objective_controls;
-- +goose StatementEnd