-- +goose Up
-- +goose StatementBegin
CREATE TABLE business_impact_assessment_analysis (
    id INT PRIMARY KEY IDENTITY(1,1),
    business_impact_assessment_id INT NOT NULL,
    outage_bucket_id INT NULL,
    impact_analysis_level_id INT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_business_impact_assessment_analysis_business_impact_assessment
        FOREIGN KEY (business_impact_assessment_id)
        REFERENCES business_impact_assessments(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT FK_business_impact_assessment_analysis_bcm_bia_outage_buckets
        FOREIGN KEY (outage_bucket_id)
        REFERENCES bcm_bia_outage_buckets(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_business_impact_assessment_analysis_bcm_impact_analysis_levels
        FOREIGN KEY (impact_analysis_level_id)
        REFERENCES bcm_impact_analysis_levels(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_business_impact_assessment_analysis_users
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE business_impact_assessment_analysis;
-- +goose StatementEnd