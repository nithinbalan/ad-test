-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_meetings (
    id INT PRIMARY KEY IDENTITY(1, 1),
    audit_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    from_time NVARCHAR(45) NULL,
    to_time NVARCHAR(45) NULL,
    other_participants NVARCHAR(MAX) NULL,
    mom NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_audit_meetings_audits FOREIGN KEY (audit_id) REFERENCES audits(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_meetings_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_meetings_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_meetings;
-- +goose StatementEnd
