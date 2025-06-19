-- +goose Up
-- +goose StatementBegin
CREATE TABLE policy_exception_risks (
    id INT PRIMARY KEY IDENTITY(1,1),
    policy_exception_id INT NOT NULL,
    risk_register_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,

    CONSTRAINT FK_policy_exception_risks_policy_exceptions 
        FOREIGN KEY (policy_exception_id) REFERENCES policy_exceptions (id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_policy_exception_risks_risk_register 
        FOREIGN KEY (risk_register_id) REFERENCES risk_register (id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_policy_exception_risks_users_created 
        FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE policy_exception_risks;
-- +goose StatementEnd
