-- +goose Up
-- +goose StatementBegin
ALTER TABLE sla_register
    ADD rto_value INT NULL,
        rto_unit_id INT NULL,
        mtpd_value INT NULL,
        mtpd_unit_id INT NULL,
        rpo_value INT NULL,
        rpo_unit_id INT NULL;

ALTER TABLE sla_register
    ADD CONSTRAINT FK_sla_register_rto_unit FOREIGN KEY (rto_unit_id) REFERENCES business_continuity_duration_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
        CONSTRAINT FK_sla_register_mtpd_unit FOREIGN KEY (mtpd_unit_id) REFERENCES business_continuity_duration_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
        CONSTRAINT FK_sla_register_rpo_unit FOREIGN KEY (rpo_unit_id) REFERENCES business_continuity_duration_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sla_register
    DROP CONSTRAINT IF EXISTS FK_sla_register_rpo_unit,
        CONSTRAINT IF EXISTS FK_sla_register_mtpd_unit,
        CONSTRAINT IF EXISTS FK_sla_register_rto_unit;

ALTER TABLE sla_register
    DROP COLUMN IF EXISTS rpo_unit_id,
                rpo_value,
                mtpd_unit_id,
                mtpd_value,
                rto_unit_id,
                rto_value;
-- +goose StatementEnd
