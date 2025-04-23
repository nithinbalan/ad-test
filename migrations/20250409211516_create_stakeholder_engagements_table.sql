-- +goose Up
-- +goose StatementBegin
CREATE TABLE stakeholder_engagements
(
    id INT PRIMARY KEY IDENTITY(1,1),
    stakeholder_id INT NOT NULL,
    title NVARCHAR(500) NOT NULL,
    engagement_strategy NVARCHAR(MAX) NULL,
    description NVARCHAR(MAX) NULL,
    review_frequency_id INT NULL,
    assigned_owner_id INT NULL,
    score INT NULL,
    evaluated_by_id INT NULL,
    evaluated_at DATE NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL,

    CONSTRAINT FK_stakeholder_engagements_stakeholders FOREIGN KEY (stakeholder_id) 
        REFERENCES stakeholders(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_stakeholder_engagements_review_frequencies FOREIGN KEY (review_frequency_id) 
        REFERENCES severity_levels(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_stakeholder_engagements_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_stakeholder_engagements_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_stakeholder_engagements_users_3 FOREIGN KEY (evaluated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stakeholder_engagements;
-- +goose StatementEnd
