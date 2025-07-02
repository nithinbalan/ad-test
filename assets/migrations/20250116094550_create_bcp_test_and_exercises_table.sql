-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_test_and_exercises (
    id INT PRIMARY KEY IDENTITY(1,1),
    business_continuity_plan_id INT NOT NULL,
    title NVARCHAR(MAX) NOT NULL,
    objective NVARCHAR(MAX) NULL,
    scenario NVARCHAR(MAX) NULL,
    review_frequency_id INT NULL,
    responsible_user_id INT NULL,
    responsible_unit_id INT NULL,
    success_criteria NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
   CONSTRAINT FK_bcp_test_and_exercises_business_continuity_plans FOREIGN KEY (business_continuity_plan_id) REFERENCES business_continuity_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_test_and_exercises_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_test_and_exercises_review_frequencies FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_test_and_exercises_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_test_and_exercises_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_test_and_exercises_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.bcp_test_and_exercises;
-- +goose StatementEnd
