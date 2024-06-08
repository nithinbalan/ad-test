-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_processes (
    id INT PRIMARY KEY IDENTITY(1, 1),
    meeting_id INT NOT NULL,
    process_id INT NOT NULL,
    created_by INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_meeting_processes_meetings FOREIGN KEY (meeting_id) REFERENCES meetings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_processes_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE meeting_processes;

-- +goose StatementEnd