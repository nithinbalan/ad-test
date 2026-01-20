-- +goose Up
-- +goose StatementBegin
CREATE TABLE consultation_request_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description_en NVARCHAR(100) NOT NULL,
    description_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(10) NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT CHK_consultation_request_statuses_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT CHK_consultation_request_statuses_type CHECK (type IN (
        'new', 
        'needs_clarification', 
        'under_processing', 
        'reply_sent', 
        'closed', 
        'reopened'
    ))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.consultation_request_statuses;
-- +goose StatementEnd
