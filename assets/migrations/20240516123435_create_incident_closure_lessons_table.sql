-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_closure_lessons
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_closure_id INT NOT NULL,
    lessons NVARCHAR(500) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_incident_closure_lessons_incident_closures FOREIGN KEY (incident_closure_id) REFERENCES incident_closures (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.incident_closure_lessons;
-- +goose StatementEnd
