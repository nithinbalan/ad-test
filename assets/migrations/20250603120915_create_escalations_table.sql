-- +goose Up
-- +goose StatementBegin
CREATE TABLE escalations (
    id INT PRIMARY KEY IDENTITY(1, 1),
    module_id INT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    status VARCHAR(10) NOT NULL,
        CONSTRAINT CHK_escalations_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_escalations_modules FOREIGN KEY (module_id) 
        REFERENCES modules(id) ON DELETE CASCADE ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.escalations;
-- +goose StatementEnd
