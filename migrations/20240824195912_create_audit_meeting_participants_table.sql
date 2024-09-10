-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_meeting_participants (
    id INT PRIMARY KEY IDENTITY(1, 1),
    audit_meeting_id INT NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT FK_audit_meeting_participants_audit_meetings FOREIGN KEY (audit_meeting_id) REFERENCES audit_meetings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_meeting_participants_users_1 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_meeting_participants;
-- +goose StatementEnd
