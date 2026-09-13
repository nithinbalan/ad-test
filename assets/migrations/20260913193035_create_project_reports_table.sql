-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_reports (
    id INT IDENTITY(1, 1) PRIMARY KEY,
    title NVARCHAR(255) NOT NULL,
    start_date DATETIME NULL,
    end_date DATETIME NULL,
    created_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_reports_users
        FOREIGN KEY (created_by_id)
        REFERENCES users (id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_reports;
-- +goose StatementEnd
