-- +goose Up
-- +goose StatementBegin
CREATE TABLE business_impact_assessment_recovery_objectives (
    id INT PRIMARY KEY IDENTITY(1,1),
    process_id INT NULL,
    asset_register_id INT NULL,
    document_id INT NULL,
    role_in_process NVARCHAR(MAX) NULL,
    rto INt NULL,
    rto_unit_id INt NULL,
    mtpd INt NULL,
    mtpd_unit_id INt NULL,
    rpo INt NULL,
    rpo_unit_id INt NULL,
    impact_description NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
   CONSTRAINT FK_business_impact_assessment_recover_objectives_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_business_impact_assessment_recover_objectives_asset_register FOREIGN KEY (asset_register_id) REFERENCES asset_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_business_impact_assessment_recover_objectives_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_business_impact_assessment_recover_objectives_time_units_1 FOREIGN KEY (rto_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_business_impact_assessment_recover_objectives_time_units_2 FOREIGN KEY (mtpd_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_business_impact_assessment_recover_objectives_time_units_3 FOREIGN KEY (rpo_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_business_impact_assessment_recover_objectives_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_business_impact_assessment_recover_objectives_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.business_impact_assessment_recovery_objectives;
-- +goose StatementEnd
