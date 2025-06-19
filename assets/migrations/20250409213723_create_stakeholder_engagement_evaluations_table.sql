-- +goose Up
-- +goose StatementBegin
CREATE TABLE stakeholder_engagement_evaluations
(
    id INT PRIMARY KEY IDENTITY(1,1),
    stakeholder_engagement_id INT NOT NULL,
    score INT NULL,
    evaluated_by_id INT NULL,
    evaluated_at DATE NULL,
    comment NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL,

    CONSTRAINT FK_stakeholder_engagement_evaluations_stakeholder_engagements FOREIGN KEY (stakeholder_engagement_id) 
        REFERENCES stakeholder_engagements(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_stakeholder_engagement_evaluations_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_stakeholder_engagement_evaluations_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_stakeholder_engagement_evaluations_users_3 FOREIGN KEY (evaluated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stakeholder_engagement_evaluations;
-- +goose StatementEnd
