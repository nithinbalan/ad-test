-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_risks (
    id INT PRIMARY KEY IDENTITY(1, 1),
    meeting_id INT NOT NULL,
    risk_register_id INT NOT NULL,
    created_by INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_meeting_risk_meetings FOREIGN KEY (meeting_id) REFERENCES meetings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_risk_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE meeting_risks;

-- +goose StatementEnd