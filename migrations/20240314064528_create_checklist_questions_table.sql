-- +goose Up
-- +goose StatementBegin
CREATE TABLE checklist_questions (
    id INT PRIMARY KEY IDENTITY(1, 1),
    checklist_id INT,
    question_library_id INT,
    title_en VARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,    
    CONSTRAINT FK_checklist_question_checklist FOREIGN KEY (checklist_id) REFERENCES checklists(id),
    CONSTRAINT FK_checklist_question_question_library FOREIGN KEY (question_library_id) REFERENCES question_library(id),

);-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.checklist_questions;
-- +goose StatementEnd
