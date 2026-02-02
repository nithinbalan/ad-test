-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_ratings (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(10) NOT NULL,
        CONSTRAINT CHK_compliance_assessment_ratings_type CHECK (type IN ('critical','high', 'medium', 'low', 'closed')),
    color_code VARCHAR(10) NOT NULL,
    description_en NVARCHAR(max) NULL,
    description_ar NVARCHAR(max) NULL,
    score_from DECIMAL(10,2) NOT NULL,
    score_to DECIMAL(10,2) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_compliance_assessment_ratings_status CHECK (status IN ('active', 'inactive'))
   
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_assessment_ratings;
-- +goose StatementEnd