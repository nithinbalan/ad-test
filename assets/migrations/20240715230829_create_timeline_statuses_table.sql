-- +goose Up
-- +goose StatementBegin
CREATE TABLE timeline_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'pending',
        CONSTRAINT CHK_timeline_statuses_type CHECK (type IN ('pending', 'completed'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.timeline_statuses;
-- +goose StatementEnd
