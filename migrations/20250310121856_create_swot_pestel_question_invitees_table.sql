-- +goose Up
-- +goose StatementBegin
CREATE TABLE swot_pestel_question_invitees (
    id INT PRIMARY KEY IDENTITY(1, 1),
    swot_pestel_question_id INT NOT NULL,
    user_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_swot_pestel_question_invitees_swot_pestel_questions FOREIGN KEY (swot_pestel_question_id) REFERENCES swot_pestel_questions(id),
    CONSTRAINT FK_swot_pestel_question_invitees_users_1 FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_swot_pestel_question_invitees_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.swot_pestel_question_invitees;
-- +goose StatementEnd
