-- +goose Up
-- +goose StatementBegin
CREATE TABLE reminders (
    id INT PRIMARY KEY IDENTITY(1, 1),
    module_id INT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(25) NOT NULL,
        CONSTRAINT CHK_reminder_type CHECK (type IN ('user-task')),
    status VARCHAR(10) NOT NULL,
        CONSTRAINT CHK_reminders_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_reminders_modules FOREIGN KEY (module_id) 
        REFERENCES modules(id) ON DELETE CASCADE ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.reminders;
-- +goose StatementEnd
