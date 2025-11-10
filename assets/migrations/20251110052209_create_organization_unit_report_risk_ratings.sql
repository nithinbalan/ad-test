-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_unit_report_risk_ratings (
    id INT PRIMARY KEY IDENTITY(1,1),
    organization_unit_report_id INT NOT NULL,
    risk_rating_title_en NVARCHAR(100) NOT NULL,
    risk_rating_title_ar NVARCHAR(100) NULL,
    risk_rating_color_code VARCHAR(9) NULL,
    count INT NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_organization_unit_report_risk_ratings_reports 
        FOREIGN KEY (organization_unit_report_id) REFERENCES organization_unit_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_organization_unit_report_risk_ratings_users_1 
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organization_unit_report_risk_ratings;
-- +goose StatementEnd
