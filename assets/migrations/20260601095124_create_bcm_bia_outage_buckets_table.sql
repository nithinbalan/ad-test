-- +goose Up
-- +goose StatementBegin

CREATE TABLE bcm_bia_outage_buckets (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    from_value INT NOT NULL,
    to_value INT NOT NULL,
    duration_unit_id INT NOT NULL,
    [order] INT NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_bcm_bia_outage_buckets_status CHECK (status IN ('active', 'inactive')),

    CONSTRAINT FK_bcm_bia_outage_buckets_business_continuity_duration_units FOREIGN KEY (duration_unit_id)
            REFERENCES business_continuity_duration_units(id)
            ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE dbo.bcm_bia_outage_buckets;

-- +goose StatementEnd
