-- +goose Up
-- +goose StatementBegin
-- Drop foreign key constraints
ALTER TABLE business_impact_assessments DROP CONSTRAINT FK_business_impact_assessments_bcm_bia_outage_buckets;
ALTER TABLE business_impact_assessments DROP CONSTRAINT FK_business_impact_assessments_bcm_impact_analysis_levels;

-- Drop new columns
ALTER TABLE business_impact_assessments DROP COLUMN outage_bucket_id;
ALTER TABLE business_impact_assessments DROP COLUMN impact_analysis_level_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Add back the columns (for rollback)
-- Add new columns
ALTER TABLE business_impact_assessments ADD impact_analysis_level_id INT NULL;
ALTER TABLE business_impact_assessments ADD outage_bucket_id INT NULL;

-- Add foreign key constraints
ALTER TABLE business_impact_assessments ADD CONSTRAINT FK_business_impact_assessments_bcm_impact_analysis_levels FOREIGN KEY (impact_analysis_level_id)
    REFERENCES bcm_impact_analysis_levels(id)
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

ALTER TABLE business_impact_assessments ADD CONSTRAINT FK_business_impact_assessments_bcm_bia_outage_buckets FOREIGN KEY (outage_bucket_id)
    REFERENCES bcm_bia_outage_buckets(id)
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;
-- +goose StatementEnd
