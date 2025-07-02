-- +goose Up
-- +goose StatementBegin
CREATE TABLE timeliness (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'ontime',
        CONSTRAINT CHK_timeliness_type CHECK (type IN ('ontime', 'delayed', 'early')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_timeliness_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.timeliness;
-- +goose StatementEnd

