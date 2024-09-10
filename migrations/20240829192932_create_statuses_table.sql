-- +goose Up
-- +goose StatementBegin
CREATE TABLE statuses
(
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    CONSTRAINT CHK_statuses_type CHECK (type IN ('pending', 'completed')),
    color_code VARCHAR(25) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_statuses_status CHECK (status IN ('active', 'inactive'))

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE statuses;
-- +goose StatementEnd
