-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    meeting_id INT NOT NULL,
    url VARCHAR(1000) NOT NULL,
    name NVARCHAR(255) NOT NULL,
    size INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_meeting_documents_meeting FOREIGN KEY (meeting_id) REFERENCES meetings(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.meeting_documents;
-- +goose StatementEnd
