-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_frequencies (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
        CONSTRAINT CHK_meeting_frequencies_type CHECK (type IN ('weekly', 'monthly', 'quarterly', 'half-yearly', 'annually')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_meeting_frequencies_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.meeting_frequencies;
-- +goose StatementEnd