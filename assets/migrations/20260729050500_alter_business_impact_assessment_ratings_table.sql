-- +goose Up
-- +goose StatementBegin
ALTER TABLE business_impact_assessment_ratings ADD overall_outage_bucket_id INT NULL;

-- Add foreign key constraints
ALTER TABLE business_impact_assessment_ratings ADD CONSTRAINT FK_business_impact_assessment_ratings_bcm_bia_outage_buckets FOREIGN KEY (overall_outage_bucket_id)
    REFERENCES bcm_bia_outage_buckets(id)
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE business_impact_assessment_ratings DROP CONSTRAINT FK_business_impact_assessment_ratings_bcm_bia_outage_buckets;
ALTER TABLE business_impact_assessment_ratings DROP COLUMN overall_outage_bucket_id;
-- +goose StatementEnd
