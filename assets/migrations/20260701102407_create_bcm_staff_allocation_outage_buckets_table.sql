-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcm_staff_allocation_outage_buckets (
    id INT IDENTITY(1,1) PRIMARY KEY,
    bcm_staff_allocation_id INT NOT NULL,
    outage_bucket_id INT NULL,
    is_skipped BIT NOT NULL DEFAULT 0,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,  

    CONSTRAINT FK_bcm_staff_allocation_outage_buckets_staff_allocation FOREIGN KEY (bcm_staff_allocation_id)
        REFERENCES bcm_staff_allocations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_bcm_staff_allocation_outage_buckets_outage_buckets FOREIGN KEY (outage_bucket_id)
        REFERENCES bcm_bia_outage_buckets(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_bcm_staff_allocation_users_users_1 FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_bcm_staff_allocation_users_users_2 FOREIGN KEY (updated_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bcm_staff_allocation_outage_buckets;
-- +goose StatementEnd
