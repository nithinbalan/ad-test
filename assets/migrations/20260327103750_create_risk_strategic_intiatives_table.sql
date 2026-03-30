-- +goose Up
-- +goose StatementBegin
CREATE TABLE risks_strategic_intiatives (
    id INT PRIMARY KEY IDENTITY(1,1),
    strategic_intiative_id INT NOT NULL,
    risk_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_risks_strategic_intiatives_strategic_intiatives
        FOREIGN KEY (strategic_intiative_id) REFERENCES strategic_intiatives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risks_strategic_intiatives_risk_register
        FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risks_strategic_intiatives_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT UQ_risks_strategic_intiatives_unique
        UNIQUE (strategic_intiative_id, risk_register_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risks_strategic_intiatives;
-- +goose StatementEnd