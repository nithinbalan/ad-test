-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_user_statuses_type CHECK (type IN ('active', 'onleave', 'resigned', 'retired', 'locked', 'deceased')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_user_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.user_statuses;
-- +goose StatementEnd

