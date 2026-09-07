-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_checklist_comments (
    id INT PRIMARY KEY IDENTITY(1,1),
    checklist_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_project_security_gate_checklist_comments_checklist
        FOREIGN KEY (checklist_id)
        REFERENCES project_security_gate_checklists(id),

    CONSTRAINT FK_project_security_gate_checklist_comments_created_by
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_checklist_comments;
-- +goose StatementEnd