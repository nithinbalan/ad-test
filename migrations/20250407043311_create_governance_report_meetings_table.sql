-- +goose Up
-- +goose StatementBegin
CREATE TABLE governance_report_meetings (
    id INT PRIMARY KEY IDENTITY(1,1),
    governance_report_id INT NOT NULL,
    meeting_title NVARCHAR(255) NOT NULL,
    meeting_date DATE NOT NULL,
    participants_count INT NOT NULL DEFAULT 0,
    key_decisions NVARCHAR(MAX) NULL,
    actions_assigned NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_governance_report_meetings_governance_reports FOREIGN KEY (governance_report_id) REFERENCES governance_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_governance_report_meetings_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE governance_report_meetings;
-- +goose StatementEnd
