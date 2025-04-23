-- +goose Up
-- +goose StatementBegin

CREATE TABLE risk_closures (
    id INT PRIMARY KEY IDENTITY(1, 1),
    risk_register_id INT NOT NULL,
    justification NVARCHAR(MAX),
    lesson_learned NVARCHAR(MAX) NULL,       
    what_worked_well NVARCHAR(MAX) NULL,     
    improvements NVARCHAR(MAX) NULL,
    takeaways NVARCHAR(MAX) NULL,        
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    updated_at DATETIME NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_risk_closures_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_closures_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_closures_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS risk_closures;

-- +goose StatementEnd
