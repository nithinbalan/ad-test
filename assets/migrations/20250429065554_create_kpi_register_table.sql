-- +goose Up
-- +goose StatementBegin
CREATE TABLE kpi_register (
    id INT PRIMARY KEY IDENTITY(1,1),
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL,
    kpi_category_id INT NOT NULL,
    responsible_user_id INT NULL,
    kpi_status_id INT NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_kpi_register_kpi_categories FOREIGN KEY (kpi_category_id) 
        REFERENCES kpi_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_register_kpi_statuses FOREIGN KEY (kpi_status_id) 
        REFERENCES kpi_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_register_users_1 FOREIGN KEY (responsible_user_id) 
        REFERENCES users(id) ON DELETE SET NULL ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_register_users_2 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_register_users_3 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.kpi_register;
-- +goose StatementEnd
