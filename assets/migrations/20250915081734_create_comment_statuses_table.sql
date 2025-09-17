-- +goose Up
-- +goose StatementBegin
CREATE TABLE comment_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL,
        CONSTRAINT CHK_comment_statuses_type CHECK (type IN ('open', 'resolved')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_comment_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.comment_statuses;
-- +goose StatementEnd
