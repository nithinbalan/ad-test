-- +goose Up
-- +goose StatementBegin
ALTER TABLE business_continuity_plans ADD parent_id INT NULL;
ALTER TABLE business_continuity_plans DROP CONSTRAINT UQ__business__321A02983C4987AB;

ALTER TABLE business_continuity_plans
ADD CONSTRAINT FK_business_continuity_plans_parent
FOREIGN KEY (parent_id) REFERENCES business_continuity_plans(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE business_continuity_plans DROP CONSTRAINT FK_business_continuity_plans_parent;
ALTER TABLE business_continuity_plans DROP COLUMN parent_id;
ADD CONSTRAINT UQ__business__321A02983C4987AB UNIQUE (reference_number);
-- +goose StatementEnd
