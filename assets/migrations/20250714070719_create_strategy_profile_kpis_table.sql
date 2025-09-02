-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategy_profile_kpis (
    id INT PRIMARY KEY IDENTITY(1,1),
    strategy_profile_id INT NOT NULL,
    strategy_profile_strategic_objective_id INT NULL,
    kpi_id INT NULL,
    weightage FLOAT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_strategy_profile_kpis_strategy_profiles FOREIGN KEY (strategy_profile_id) 
        REFERENCES strategy_profiles(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_kpis_strategy_profile_strategic_objectives FOREIGN KEY (strategy_profile_strategic_objective_id) 
        REFERENCES strategy_profile_strategic_objectives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_kpis_kpi_register FOREIGN KEY (kpi_id) 
        REFERENCES kpi_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_kpis_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_kpis_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.strategy_profile_kpis;
-- +goose StatementEnd
