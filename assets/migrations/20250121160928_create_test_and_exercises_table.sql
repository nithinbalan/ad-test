-- +goose Up
-- +goose StatementBegin
CREATE TABLE test_and_exercises (
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    business_continuity_plan_id INT NOT NULL,
    bcp_test_and_exercise_id INT NOT NULL,
    responsible_user_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    testing_date DATE NOT NULL,
    description NVARCHAR(MAX) NULL,
    test_and_exercise_status_id INT NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_test_and_exercises_business_continuity_plans FOREIGN KEY (business_continuity_plan_id) REFERENCES business_continuity_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercises_bcp_test_and_exercises FOREIGN KEY (bcp_test_and_exercise_id) REFERENCES bcp_test_and_exercises(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercises_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercises_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercises_test_and_exercise_statuses FOREIGN KEY (test_and_exercise_status_id) REFERENCES test_and_exercise_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercises_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercises_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.test_and_exercises;
-- +goose StatementEnd
