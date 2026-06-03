-- +goose Up
-- +goose StatementBegin

CREATE TABLE bcm_bia_impact_levels (
    id INT PRIMARY KEY IDENTITY(1, 1),
    impact_analysis_level_id INT NOT NULL,
    classification_en VARCHAR(50) NOT NULL,
    classification_ar NVARCHAR(50) NOT NULL,
    is_critical BIT NOT NULL DEFAULT 0,
    rto_from_value INT NULL,
    rto_to_value INT NULL,
    rto_condition VARCHAR(20) NULL,
    rto_unit_id INT NULL,
    mtpd_from_value INT NULL,
    mtpd_to_value INT NULL,
    mtpd_condition VARCHAR(20) NULL,
    mtpd_unit_id INT NULL,
    rpo_from_value INT NULL,
    rpo_to_value INT NULL,
    rpo_condition VARCHAR(20) NULL,
    rpo_unit_id INT NULL,
    rsp_from_value INT NULL,
    rsp_to_value INT NULL,
    rsp_condition VARCHAR(20) NULL,
    rsp_unit_id INT NULL,
    [order] INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_bcm_bia_impact_levels_status CHECK (status IN ('active', 'inactive')),
        
    CONSTRAINT FK_bcm_bia_impact_levels_bcm_impact_analysis_levels FOREIGN KEY (impact_analysis_level_id)
            REFERENCES bcm_impact_analysis_levels(id)
             ON DELETE NO ACTION ON UPDATE NO ACTION,
    
    CONSTRAINT FK_bcm_bia_impact_levels_business_continuity_duration_units_rto FOREIGN KEY (rto_unit_id)
            REFERENCES business_continuity_duration_units(id)
             ON DELETE NO ACTION ON UPDATE NO ACTION,
    
    CONSTRAINT FK_bcm_bia_impact_levels_business_continuity_duration_units_mtpd FOREIGN KEY (mtpd_unit_id)
            REFERENCES business_continuity_duration_units(id)
             ON DELETE NO ACTION ON UPDATE NO ACTION,
    
    CONSTRAINT FK_bcm_bia_impact_levels_business_continuity_duration_units_rpo FOREIGN KEY (rpo_unit_id)
            REFERENCES business_continuity_duration_units(id)
             ON DELETE NO ACTION ON UPDATE NO ACTION,
    
    CONSTRAINT FK_bcm_bia_impact_levels_business_continuity_duration_units_rsp FOREIGN KEY (rsp_unit_id)
            REFERENCES business_continuity_duration_units(id)
             ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE dbo.bcm_bia_impact_levels;

-- +goose StatementEnd
