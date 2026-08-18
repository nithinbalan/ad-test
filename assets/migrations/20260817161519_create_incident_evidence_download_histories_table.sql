-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_evidence_download_histories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_evidence_id INT NOT NULL,
    downloaded_at DATETIME NOT NULL,
    downloaded_by INT NOT NULL,

    CONSTRAINT FK_incident_evidence_download_histories_incident_evidences
        FOREIGN KEY (incident_evidence_id) REFERENCES incident_evidences(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_evidence_download_histories_users
        FOREIGN KEY (downloaded_by) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_evidence_download_histories;
-- +goose StatementEnd
