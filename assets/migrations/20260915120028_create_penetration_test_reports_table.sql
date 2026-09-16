-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_reports (
    id INT IDENTITY(1,1) PRIMARY KEY,

    title NVARCHAR(255) NOT NULL,
    start_date DATE NULL,
    end_date DATE NULL,

    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT FK_penetration_test_reports_users_created
        FOREIGN KEY (created_by_id)
        REFERENCES users (id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE penetration_test_reports;
-- +goose StatementEnd
