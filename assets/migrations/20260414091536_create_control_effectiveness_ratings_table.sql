-- +goose Up
-- +goose StatementBegin
CREATE TABLE control_effectiveness_ratings (
    id INT IDENTITY PRIMARY KEY,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
     type VARCHAR(100) NOT NULL,
        CONSTRAINT CHK_control_effectiveness_ratings_type 
        CHECK (type IN ('control_missing', 'ineffective', 'improvement_needed', 'partially_effective', 'very_effective')),
    min_score DECIMAL(5,2) NOT NULL,
    max_score DECIMAL(5,2) NOT NULL,
    reduction_percent DECIMAL(5,1) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS control_effectiveness_ratings;
-- +goose StatementEnd
