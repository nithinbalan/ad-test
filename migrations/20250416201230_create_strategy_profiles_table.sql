-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategy_profiles (
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(45) NULL,
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL,
    vision NVARCHAR(MAX) NULL,
    mission NVARCHAR(MAX) NULL,
    responsible_user_id INT NULL,
    start_date DATE NULL,
    end_date DATE NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_strategy_profiles_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profiles_users_2 FOREIGN KEY (responsible_user_id) 
        REFERENCES users(id) ON DELETE SET NULL ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profiles_users_3 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.strategy_profiles;
-- +goose StatementEnd
