-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_strategic_objectives  (
    id INT PRIMARY KEY IDENTITY(1,1),
    strategic_objective_id INT NOT NULL,
    meeting_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    CONSTRAINT FK_meeting_strategic_objectives_strategic_objective FOREIGN KEY (strategic_objective_id) REFERENCES strategic_objectives(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_strategic_objectives_meetings FOREIGN KEY (meeting_id) REFERENCES meetings(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_strategic_objectives_user_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE meeting_strategic_objectives;
-- +goose StatementEnd
