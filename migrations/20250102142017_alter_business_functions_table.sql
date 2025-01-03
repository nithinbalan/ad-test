-- +goose Up
-- +goose StatementBegin
ALTER TABLE business_functions ADD review_frequency_id INT NULL;
ALTER TABLE business_functions ADD last_review_date DATE NULL;
ALTER TABLE business_functions ADD next_review_date DATE NULL;
ALTER TABLE business_functions ADD bcp_status TINYINT NOT NULL DEFAULT 0;


ALTER TABLE business_functions ADD CONSTRAINT FK_business_functions_review_frequencies 
    FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies(id) ON DELETE SET NULL ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE business_functions DROP CONSTRAINT FK_business_functions_review_frequencies;
ALTER TABLE business_functions DROP COLUMN review_frequency_id;
ALTER TABLE business_functions DROP COLUMN last_review_date;
ALTER TABLE business_functions DROP COLUMN next_review_date;
ALTER TABLE business_functions DROP COLUMN bcp_status;
-- +goose StatementEnd
