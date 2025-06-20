-- +goose Up
-- +goose StatementBegin

-- Add new columns
ALTER TABLE strategic_objectives
ADD strategic_objective_category_id INT NULL,
    next_review_date DATE NULL,
    review_status_id INT NULL;

-- Add foreign keys
ALTER TABLE strategic_objectives
ADD CONSTRAINT FK_strategic_objectives_categories FOREIGN KEY (strategic_objective_category_id)
    REFERENCES strategic_objective_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategic_objectives_review_status FOREIGN KEY (review_status_id)
    REFERENCES review_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop foreign keys first
ALTER TABLE strategic_objectives DROP CONSTRAINT FK_strategic_objectives_categories;
ALTER TABLE strategic_objectives DROP CONSTRAINT FK_strategic_objectives_review_status;

-- Drop columns
ALTER TABLE strategic_objectives 
DROP COLUMN strategic_objective_category_id,
            next_review_date,
            review_status_id;

-- +goose StatementEnd
