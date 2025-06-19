-- +goose Up
-- +goose StatementBegin
CREATE TABLE third_party_assessments
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    title NVARCHAR(500) NOT NULL,
    third_party_register_id INT NOT NULL,
    question_set_id INT NOT NULL,
    message NVARCHAR(MAX),
    due_date DATE,
    third_party_assessment_status_id INT NOT NULL,

    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_third_party_assessments_third_party_register FOREIGN KEY (third_party_register_id) REFERENCES third_party_register (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_assessments_question_sets FOREIGN KEY (question_set_id) REFERENCES question_sets (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_assessments_statuses FOREIGN KEY (third_party_assessment_status_id) REFERENCES third_party_assessment_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_assessments_user1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_assessments_user2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE third_party_assessments;
-- +goose StatementEnd
