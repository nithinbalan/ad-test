-- +goose Up
-- +goose StatementBegin
CREATE TABLE policy_reports (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_policy_reports_documents FOREIGN KEY (document_id) REFERENCES documents (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_policy_reports_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE policy_reports;
-- +goose StatementEnd
