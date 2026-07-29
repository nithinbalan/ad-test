-- +goose Up
-- +goose StatementBegin
ALTER TABLE bcm_bia_impact_levels ADD outage_bucket_id INT NULL;

-- Add foreign key constraints
ALTER TABLE bcm_bia_impact_levels ADD CONSTRAINT FK_bcm_bia_impact_levels_bcm_bia_outage_buckets FOREIGN KEY (outage_bucket_id)
    REFERENCES bcm_bia_outage_buckets(id)
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE bcm_bia_impact_levels DROP CONSTRAINT FK_bcm_bia_impact_levels_bcm_bia_outage_buckets;
ALTER TABLE bcm_bia_impact_levels DROP COLUMN outage_bucket_id;
-- +goose StatementEnd
