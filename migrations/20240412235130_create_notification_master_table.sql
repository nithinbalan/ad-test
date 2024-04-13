-- +goose Up
-- +goose StatementBegin
CREATE TABLE notification_masters (
    id INT PRIMARY KEY IDENTITY(1,1),
    modules_id INT,
    identifier VARCHAR(100),
    title_en VARCHAR(100),
    title_ar NVARCHAR(100),
    CONSTRAINT FK_notification_masters_modules_id FOREIGN KEY (modules_id) REFERENCES modules(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_masters;
-- +goose StatementEnd
