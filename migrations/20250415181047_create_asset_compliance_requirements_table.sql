-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_compliance_requirements (
    id INT PRIMARY KEY IDENTITY(1, 1),
    asset_register_id INT NOT NULL,
    compliance_clause_id INT NULL,
    compliance_requirement_id INT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_asset_compliance_requirements_asset_register FOREIGN KEY (asset_register_id) REFERENCES asset_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_compliance_requirements_compliance_clause FOREIGN KEY (compliance_clause_id) REFERENCES compliance_clauses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_compliance_requirements_compliance_requirement FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_compliance_requirements_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE asset_compliance_requirements;
-- +goose StatementEnd
