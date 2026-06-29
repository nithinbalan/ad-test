-- +goose Up
-- +goose StatementBegin
ALTER TABLE business_impact_assessment_ratings
    ADD is_edited BIT NOT NULL CONSTRAINT DF_business_impact_assessment_ratings_is_edited DEFAULT 0,
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
        updated_by_id INT NULL,
        updated_at DATETIME NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE business_impact_assessment_ratings ADD
    CONSTRAINT FK_business_impact_assessment_ratings_duration_units_rto FOREIGN KEY (rto_unit_id) REFERENCES business_continuity_duration_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_business_impact_assessment_ratings_duration_units_mtpd FOREIGN KEY (mtpd_unit_id) REFERENCES business_continuity_duration_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_business_impact_assessment_ratings_duration_units_rpo FOREIGN KEY (rpo_unit_id) REFERENCES business_continuity_duration_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_business_impact_assessment_ratings_duration_units_rsp FOREIGN KEY (rsp_unit_id) REFERENCES business_continuity_duration_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_business_impact_assessment_ratings_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE business_impact_assessment_ratings DROP CONSTRAINT
    FK_business_impact_assessment_ratings_updated_by,
    FK_business_impact_assessment_ratings_duration_units_rsp,
    FK_business_impact_assessment_ratings_duration_units_rpo,
    FK_business_impact_assessment_ratings_duration_units_mtpd,
    FK_business_impact_assessment_ratings_duration_units_rto,
    DF_business_impact_assessment_ratings_is_edited;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE business_impact_assessment_ratings DROP COLUMN
    updated_at,
    updated_by_id,
    rsp_unit_id,
    rsp_condition,
    rsp_to_value,
    rsp_from_value,
    rpo_unit_id,
    rpo_condition,
    rpo_to_value,
    rpo_from_value,
    mtpd_unit_id,
    mtpd_condition,
    mtpd_to_value,
    mtpd_from_value,
    rto_unit_id,
    rto_condition,
    rto_to_value,
    rto_from_value,
    is_edited;
-- +goose StatementEnd
