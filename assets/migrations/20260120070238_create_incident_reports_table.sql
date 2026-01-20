-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_reports (
    id INT IDENTITY(1,1) PRIMARY KEY,

    title NVARCHAR(255) NOT NULL,
    start_date DATE NULL,
    end_date DATE NULL,

    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),


    CONSTRAINT FK_incident_reports_users_created
        FOREIGN KEY (created_by_id)
        REFERENCES users (id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_reports;
-- +goose StatementEnd
