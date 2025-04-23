-- +goose Up
-- +goose StatementBegin
CREATE TABLE stakeholder_needs_and_expectation_evaluations (
    id INT PRIMARY KEY IDENTITY(1, 1),
    stakeholder_needs_and_expectation_id INT NOT NULL,
    evaluation_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    score INT NULL,
    evaluated_by_id INT NULL,
    evaluated_at DATE NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_stakeholder_needs_and_expectation_evaluations_needs FOREIGN KEY (stakeholder_needs_and_expectation_id) REFERENCES stakeholder_needs_and_expectations (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_stakeholder_needs_and_expectation_evaluations_evaluation_statuses FOREIGN KEY (evaluation_status_id) REFERENCES evaluation_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_stakeholder_needs_and_expectation_evaluations_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_stakeholder_needs_and_expectation_evaluations_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_stakeholder_needs_and_expectation_evaluations_users_3 FOREIGN KEY (evaluated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stakeholder_needs_and_expectation_evaluations;
-- +goose StatementEnd
