-- +goose Up
-- +goose StatementBegin

CREATE TABLE business_impact_assessment_ratings (
    id INT PRIMARY KEY IDENTITY(1, 1),
    process_id INT NOT NULL,
    overall_impact_level_id INT NOT NULL,
    bcm_bia_impact_level_id INT NOT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_business_impact_assessment_ratings_processes FOREIGN KEY (process_id) REFERENCES processes (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_business_impact_assessment_ratings_overall_impact_levels FOREIGN KEY (overall_impact_level_id) REFERENCES bcm_impact_analysis_levels (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_business_impact_assessment_ratings_bcm_bia_impact_levels FOREIGN KEY (bcm_bia_impact_level_id) REFERENCES bcm_bia_impact_levels (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_business_impact_assessment_ratings_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
    
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE business_impact_assessment_ratings;
-- +goose StatementEnd
