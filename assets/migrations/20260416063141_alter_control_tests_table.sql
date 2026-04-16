-- +goose Up
-- +goose StatementBegin

ALTER TABLE control_tests
ADD 
    control_effectiveness_score DECIMAL(5,2) NULL,
    control_effectiveness_rating_id INT NULL,
    risk_reduction DECIMAL(5,2) NULL;

ALTER TABLE control_tests
ADD CONSTRAINT FK_control_tests_effectiveness_rating_id 
FOREIGN KEY (control_effectiveness_rating_id) REFERENCES control_effectiveness_ratings(id);
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

-- Drop FK first
ALTER TABLE control_tests
DROP CONSTRAINT FK_control_tests_effectiveness_rating_id;

-- Then drop columns
ALTER TABLE control_tests
DROP COLUMN 
    control_effectiveness_score,
    control_effectiveness_rating_id,
    risk_reduction;

-- +goose StatementEnd