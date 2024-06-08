-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_controls (
    id INT PRIMARY KEY IDENTITY(1, 1),
    meeting_id INT NOT NULL,
    control_id INT NOT NULL,
    created_by INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_meeting_controls_meetings FOREIGN KEY (meeting_id) REFERENCES meetings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_controls_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE meeting_controls;

-- +goose StatementEnd