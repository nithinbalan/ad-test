-- +goose Up
-- +goose StatementBegin

-- Drop old foreign key constraint
ALTER TABLE kri_register
DROP CONSTRAINT FK_kri_register_risk_sub_categories;

-- Drop old column
ALTER TABLE kri_register
DROP COLUMN risk_sub_category_id;

-- Add new column
ALTER TABLE kri_register
ADD kri_sub_category_id INT NULL;

-- Add new foreign key constraint
ALTER TABLE kri_register
ADD CONSTRAINT FK_kri_register_kri_sub_categories
FOREIGN KEY (kri_sub_category_id)
REFERENCES kri_sub_categories(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Drop new foreign key constraint
ALTER TABLE kri_register
DROP CONSTRAINT FK_kri_register_kri_sub_categories;

-- Drop new column
ALTER TABLE kri_register
DROP COLUMN kri_sub_category_id;

-- Add old column back
ALTER TABLE kri_register
ADD risk_sub_category_id INT NULL;

-- Add old foreign key constraint back
ALTER TABLE kri_register
ADD CONSTRAINT FK_kri_register_risk_sub_categories
FOREIGN KEY (risk_sub_category_id)
REFERENCES risk_sub_categories(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- +goose StatementEnd