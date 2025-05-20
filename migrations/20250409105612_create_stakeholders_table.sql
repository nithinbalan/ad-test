-- +goose Up
-- +goose StatementBegin
CREATE TABLE stakeholders (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL,

    CONSTRAINT FK_stakeholders_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_stakeholders_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.stakeholders;
-- +goose StatementEnd
