-- +goose Up
-- +goose StatementBegin
CREATE TABLE checklist_questions (
    id INT PRIMARY KEY IDENTITY(1, 1),
    checklist_id INT,
    question_group_id INT,
    question_en NVARCHAR(500) NOT NULL,
    question_ar NVARCHAR(500) NOT NULL,    
    CONSTRAINT FK_checklist_questions_checklists FOREIGN KEY (checklist_id) REFERENCES checklists(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_checklist_questions_question_groups FOREIGN KEY (question_group_id) REFERENCES question_groups(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.checklist_questions;
-- +goose StatementEnd
