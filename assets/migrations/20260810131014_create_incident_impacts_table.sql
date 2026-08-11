-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_impacts (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(MAX) NOT NULL,
    title_ar NVARCHAR(MAX) NOT NULL,
    description NVARCHAR(MAX) NULL,
    type VARCHAR(50) NOT NULL,
        CONSTRAINT CHK_incident_impacts_type CHECK (type IN ('confidentiality', 'integrity', 'availability')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_incident_impacts_status CHECK (status IN ('active', 'inactive')),
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_incident_impacts_created_by
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_impacts;
-- +goose StatementEnd
