-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_timelines (
    id INT PRIMARY KEY IDENTITY(1, 1),
    audit_id INT NOT NULL,
    audit_timeline_type_id INT NOT NULL,
    timeline_status_id INT NOT NULL,
    performed_by_id INT NULL,
    performed_at DATETIME NULL,
    CONSTRAINT FK_audit_timelines_audits FOREIGN KEY (audit_id) REFERENCES audits(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_timelines_audit_timeline_types FOREIGN KEY (audit_timeline_type_id) REFERENCES audit_timeline_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_timelines_timeline_statuses FOREIGN KEY (timeline_status_id) REFERENCES timeline_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_timelines_users FOREIGN KEY (performed_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_timelines;
-- +goose StatementEnd
