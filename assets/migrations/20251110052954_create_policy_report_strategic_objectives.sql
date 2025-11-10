-- +goose Up
-- +goose StatementBegin
CREATE TABLE policy_report_strategic_objectives (
    id INT PRIMARY KEY IDENTITY(1,1),
    policy_report_id INT NOT NULL,
    strategic_objective_title_en NVARCHAR(255) NOT NULL,
    strategic_objective_title_ar NVARCHAR(255) NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    start_date VARCHAR(50) NULL,
    end_date VARCHAR(50) NULL,
    strategic_objective_category_title_en NVARCHAR(255) NULL,
    strategic_objective_category_title_ar NVARCHAR(255) NULL,
    responsible_user_name_en NVARCHAR(255) NULL,
    responsible_user_name_ar NVARCHAR(255) NULL,
    responsible_unit_title_en NVARCHAR(255) NULL,
    responsible_unit_title_ar NVARCHAR(255) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_policy_report_strategic_objectives_reports
        FOREIGN KEY (policy_report_id) REFERENCES policy_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_policy_report_strategic_objectives_users_1
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE policy_report_strategic_objectives;
-- +goose StatementEnd
