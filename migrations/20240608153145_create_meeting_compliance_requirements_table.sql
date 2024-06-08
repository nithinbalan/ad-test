-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_compliance_requirements
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    meeting_id INT NOT NULL,
    compliance_requirement_id INT NOT NULL,
    created_by INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_meeting_compliance_requirements_meetings FOREIGN KEY (meeting_id) REFERENCES meetings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_compliance_requirements_compliance_requirements FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE meeting_compliance_requirements;

-- +goose StatementEnd