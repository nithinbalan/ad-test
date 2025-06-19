-- +goose Up
-- +goose StatementBegin
CREATE TABLE business_impact_assessments (
    id INT PRIMARY KEY IDENTITY(1,1),
    process_id INT NULL,
    asset_register_id INT NULL,
    impact_category_id INT NOT NULL,
    risk_level_id INT NOT NULL,
    cost INT NULL,
    comments NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
   CONSTRAINT FK_business_impact_assessments_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_business_impact_assessments_asset_register FOREIGN KEY (asset_register_id) REFERENCES asset_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_business_impact_assessments_impact_category FOREIGN KEY (impact_category_id) REFERENCES impact_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_business_impact_assessments_risk_level FOREIGN KEY (risk_level_id) REFERENCES risk_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_business_impact_assessments_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_business_impact_assessments_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.business_impact_assessments;
-- +goose StatementEnd
