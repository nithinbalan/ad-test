-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_check_answers (
    id INT PRIMARY KEY IDENTITY(1, 1),
    sla_check_question_id INT NOT NULL,
    compliance_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_sla_check_answers_sla_check_questions FOREIGN KEY (sla_check_question_id) REFERENCES sla_check_questions (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_check_answers_compliance_statuses FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_check_answers_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_check_answers_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sla_check_answers;
-- +goose StatementEnd
