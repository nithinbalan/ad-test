-- +goose Up
-- +goose StatementBegin

ALTER TABLE strategic_objectives
ADD reference_number VARCHAR(45) NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    start_date DATE NULL,
    end_date DATE NULL,
    responsible_user_id INT NULL,
    responsible_unit_id INT NULL,
    review_user_id INT NULL,
    review_frequency_id INT NULL;

-- Add foreign key constraints
ALTER TABLE strategic_objectives
ADD CONSTRAINT FK_strategic_objectives_responsible_user FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategic_objectives_responsible_unit FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategic_objectives_review_user FOREIGN KEY (review_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategic_objectives_review_frequency FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop foreign keys first
ALTER TABLE strategic_objectives DROP CONSTRAINT FK_strategic_objectives_responsible_user;
ALTER TABLE strategic_objectives DROP CONSTRAINT FK_strategic_objectives_responsible_unit;
ALTER TABLE strategic_objectives DROP CONSTRAINT FK_strategic_objectives_review_user;
ALTER TABLE strategic_objectives DROP CONSTRAINT FK_strategic_objectives_review_frequency;

-- Drop added columns
ALTER TABLE strategic_objectives
DROP COLUMN reference_number,
             description_en,
             description_ar,
             start_date,
             end_date,
             responsible_user_id,
             responsible_unit_id,
             review_user_id,
             review_frequency_id;

-- +goose StatementEnd
