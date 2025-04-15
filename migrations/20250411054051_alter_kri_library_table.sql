-- +goose Up
-- +goose StatementBegin
ALTER TABLE kri_library
    DROP CONSTRAINT FK_kri_library_kri_categories,
    DROP CONSTRAINT FK_kri_library_review_frequencies,
    DROP COLUMN kri_category_id,
    DROP COLUMN review_frequency_id,
    DROP COLUMN calculation_type,
    DROP COLUMN source;
 
ALTER TABLE kri_library
    ADD kri_category NVARCHAR(255) NULL,
    ADD kri_unit NVARCHAR(45) NULL;
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
ALTER TABLE kri_library
    DROP COLUMN kri_category,
    DROP COLUMN kri_unit;
 
ALTER TABLE kri_library
    ADD kri_category_id INT NOT NULL,
    ADD review_frequency_id INT NULL,
    ADD calculation_type VARCHAR(10) NOT NULL DEFAULT 'automatic',
    ADD source VARCHAR(45) NULL;
 
ALTER TABLE kri_library
    ADD CONSTRAINT CHK_kri_library_calculation_type CHECK (calculation_type IN ('automatic', 'manual')),
    ADD CONSTRAINT FK_kri_library_kri_categories FOREIGN KEY (kri_category_id) REFERENCES kri_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    ADD CONSTRAINT FK_kri_library_review_frequencies FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd